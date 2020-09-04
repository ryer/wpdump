package wpdump

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	report    Report
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

func (dumper *WPDumper) SetReport(report Report) {
	dumper.report = report
}

func (dumper *WPDumper) Dump(path Path) ([]string, error) {
	const RetryCount = 2

	const RetryWaitTime = 5 * time.Second

	dumper.client.
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime)

	files := make([]string, 0, 1000)

	for page := 1; ; page++ {
		url := fmt.Sprintf("%v/%v", dumper.baseURL, path)

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

		response, err := request.Get(url)
		if err != nil {
			return nil, err
		}

		if response.StatusCode() != http.StatusOK {
			return nil, ErrNotOK
		}

		total, err := strconv.Atoi(response.Header().Get("X-WP-TotalPages"))
		if err != nil {
			return nil, ErrNoWPTotalPages
		}

		body := response.Body()
		filename := fmt.Sprintf("%v/%v%04d.json", dumper.outputDir, path, page)

		err = ioutil.WriteFile(filename, body, 0o644)
		if err != nil {
			return nil, err
		}

		if dumper.report != nil {
			dumper.report(path, filename)
		}

		files = append(files, filename)

		if total <= page {
			break
		}
	}

	return files, nil
}
