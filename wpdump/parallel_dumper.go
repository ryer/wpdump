package wpdump

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type ParallelDumper struct {
	dumper   *WPDumper
	reporter Reporter
	parallel int
}

func NewParallelDumper(dumper *WPDumper, parallel int) *ParallelDumper {
	return &ParallelDumper{
		dumper:   dumper,
		parallel: parallel,
	}
}

func (pd *ParallelDumper) SetReporter(reporter Reporter) {
	pd.reporter = reporter
	pd.dumper.SetReporter(reporter) // Pass through to the inner dumper
}

func (pd *ParallelDumper) getTotalPages(path Path) (int, error) {
	request := pd.dumper.buildRequest(1) // Get page 1 to check headers
	url := fmt.Sprintf("%v/%v", pd.dumper.baseURL, path)
	response, err := request.Get(url)
	if err != nil {
		return 0, err
	}
	if response.StatusCode() != http.StatusOK {
		return 0, ErrNotOK
	}

	return strconv.Atoi(response.Header().Get("X-Wp-Totalpages"))
}

func (pd *ParallelDumper) downloadPage(ctx context.Context, path Path, page int) (string, error) {
	request := pd.dumper.buildRequest(page).SetContext(ctx)
	url := fmt.Sprintf("%v/%v", pd.dumper.baseURL, path)
	response, err := request.Get(url)
	if err != nil {
		return "", err
	}

	_, filename, err := pd.dumper.processResponse(response, path, page)

	return filename, err
}

func (pd *ParallelDumper) worker(ctx context.Context, cancel context.CancelFunc, path Path, jobs <-chan int, results chan<- string, errs chan<- error) {
	var once sync.Once
	for {
		select {
		case <-ctx.Done():
			return
		case page, ok := <-jobs:
			if !ok {
				return
			}
			filename, err := pd.downloadPage(ctx, path, page)
			if err != nil {
				errs <- err
				once.Do(cancel)

				return
			}
			if pd.reporter != nil {
				pd.reporter.Success(path, filename)
			}
			results <- filename
		}
	}
}

func (pd *ParallelDumper) reportError(err error) error {
	if pd.reporter != nil {
		pd.reporter.Error(err)
	}

	return err
}

func (pd *ParallelDumper) Dump(path Path) ([]string, error) {
	totalPages, err := pd.getTotalPages(path)
	if err != nil {
		return nil, pd.reportError(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int, totalPages)
	results := make(chan string, totalPages)
	errs := make(chan error, pd.parallel)

	var wg sync.WaitGroup
	for i := 0; i < pd.parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pd.worker(ctx, cancel, path, jobs, results, errs)
		}()
	}

	for i := 1; i <= totalPages; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)
	close(errs)

	if err := <-errs; err != nil {
		return nil, pd.reportError(err)
	}

	files := make([]string, 0, totalPages)
	for filename := range results {
		files = append(files, filename)
	}
	sort.Strings(files)

	return files, nil
}
