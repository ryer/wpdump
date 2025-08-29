package wpdump

import (
	"fmt"
	"os"
	"strings"
)

type WPMergeDumper struct {
	dumper    IDumper
	reporter  Reporter
	outputDir string
}

func NewMergeDumper(dumper IDumper, outputDir string) *WPMergeDumper {
	return &WPMergeDumper{
		dumper:    dumper,
		outputDir: outputDir,
	}
}

func (merger *WPMergeDumper) SetReporter(reporter Reporter) {
	merger.reporter = reporter
	merger.dumper.SetReporter(reporter)
}

func (merger *WPMergeDumper) Dump(path Path) ([]string, error) {
	return merger.merge(merger.dumper.Dump, path)
}

func (merger *WPMergeDumper) mergeFiles(files []string, outFile *os.File) error {
	first := true
	for _, name := range files {
		data, err := os.ReadFile(name)
		if err != nil {
			return err
		}

		inner := strings.Trim(string(data), "[] \n\r") // Trim whitespace as well
		if inner != "" {
			if !first {
				if _, err := outFile.WriteString(",\n"); err != nil {
					return err
				}
			}
			if _, err := outFile.WriteString(inner); err != nil {
				return err
			}
			first = false
		}
	}

	return nil
}

func (merger *WPMergeDumper) reportError(err error) error {
	if merger.reporter != nil {
		merger.reporter.Error(err)
	}

	return err
}

func (merger *WPMergeDumper) merge(dump func(path Path) ([]string, error), path Path) ([]string, error) {
	files, err := dump(path)
	if err != nil {
		return nil, merger.reportError(err)
	}

	filename := fmt.Sprintf("%v/%v.json", merger.outputDir, path)
	outFile, err := os.Create(filename)
	if err != nil {
		return nil, merger.reportError(err)
	}
	defer outFile.Close()

	if _, err = outFile.WriteString("["); err != nil {
		return nil, merger.reportError(err)
	}

	if err := merger.mergeFiles(files, outFile); err != nil {
		return nil, merger.reportError(err)
	}

	if _, err = outFile.WriteString("]"); err != nil {
		return nil, merger.reportError(err)
	}

	for _, name := range files {
		if err := os.Remove(name); err != nil {
			if merger.reporter != nil {
				merger.reporter.Warn("failed to remove temporary file: " + name)
			}
		}
	}

	if merger.reporter != nil {
		merger.reporter.Success(path, filename)
	}

	return []string{filename}, nil
}
