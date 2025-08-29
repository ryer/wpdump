package wpdump

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	ErrNotOK          = errors.New("HTTP Status is not OK")
	ErrNoWPTotalPages = errors.New("failed to retrieve 'X-WP-TotalPages'")
)

type WPDumper struct {
	baseURL   string
	outputDir string
	reporter  Reporter
	client    *resty.Client
	embed     bool
}

func NewDumper(baseURL string, outputDir string, embed bool) *WPDumper {
	return &WPDumper{
		baseURL:   baseURL,
		outputDir: outputDir,
		embed:     embed,
		client:    resty.New(),
	}
}

func (dumper *WPDumper) SetReporter(reporter Reporter) {
	dumper.reporter = reporter
}

func (dumper *WPDumper) buildRequest(page int) *resty.Request {
	request := dumper.client.R()
	request.SetQueryParams(map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": "100",
		"orderby":  "id",
		"order":    "asc",
		"xrandom":  strconv.FormatInt(time.Now().Unix(), 36),
	})

	if dumper.embed {
		request.SetQueryParam("_embed", "1")
	}

	return request
}

func (dumper *WPDumper) processResponse(response *resty.Response, path Path, page int) (int, string, error) {
	if response.StatusCode() != http.StatusOK {
		return 0, "", ErrNotOK
	}

	total, err := strconv.Atoi(response.Header().Get("X-WP-TotalPages"))
	if err != nil {
		return 0, "", ErrNoWPTotalPages
	}

	body := response.Body()
	filename := fmt.Sprintf("%v/%v%04d.json", dumper.outputDir, path, page)

	err = os.WriteFile(filename, body, 0o644)
	if err != nil {
		return 0, "", err
	}

	return total, filename, nil
}

func (dumper *WPDumper) reportError(err error) error {
	if dumper.reporter != nil {
		dumper.reporter.Error(err)
	}

	return err
}

func (dumper *WPDumper) Dump(path Path) ([]string, error) {
	const RetryCount = 2
	const RetryWaitTime = 5 * time.Second

	dumper.client.
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime)

	files := make([]string, 0, 1000)
	url := fmt.Sprintf("%v/%v", dumper.baseURL, path)

	for page := 1; ; page++ {
		request := dumper.buildRequest(page)
		response, err := request.Get(url)
		if err != nil {
			return nil, dumper.reportError(err)
		}

		total, filename, err := dumper.processResponse(response, path, page)
		if err != nil {
			return nil, dumper.reportError(err)
		}

		if dumper.reporter != nil {
			dumper.reporter.Success(path, filename)
		}

		files = append(files, filename)

		if total <= page {
			break
		}
	}

	return files, nil
}
