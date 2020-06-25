package wpdump

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type WPMergeDumper struct {
	dumper    *WPDumper
	report    Report
	outputDir string
}

func NewMergeDumper(dumper *WPDumper, outputDir string) *WPMergeDumper {
	return &WPMergeDumper{
		dumper:    dumper,
		outputDir: outputDir,
	}
}

func (merger *WPMergeDumper) SetReport(report Report) {
	merger.report = report
	merger.dumper.SetReport(nil)
}

func (merger *WPMergeDumper) Dump(path Path) ([]string, error) {
	return merger.merge(merger.dumper.Dump, path)
}

func (merger *WPMergeDumper) merge(dump func(path Path) ([]string, error), path Path) ([]string, error) {
	err := exec.Command("jq", "--help").Run()
	if err != nil {
		return nil, err
	}

	files, err := dump(path)
	if err != nil {
		return nil, err
	}

	args := append([]string{"-c", "-s", "[.[][]]"}, files...)
	cmd := exec.Command("jq", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%v/%v.json", merger.outputDir, path)
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		return nil, err
	}

	for _, name := range files {
		err := os.Remove(name)
		if err != nil {
			return nil, err
		}
	}

	if merger.report != nil {
		merger.report(path, filename)
	}

	return []string{filename}, nil
}
