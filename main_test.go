package main

import (
	"github.com/ryer/wpdump/wpdump"
	"reflect"
	"testing"
)

func TestDecideDumpTarget(t *testing.T) {
	flags := &appFlags{}

	flags.categories = true
	a := decideDumpTarget(flags)
	if len(a) != 1 {
		t.Fatalf("did not select a 1 path (%v)", len(a))
	}
	if a[0] != wpdump.Categories {
		t.Fatalf("did not select categories (%v)", a[0])
	}

	flags.all = true
	a = decideDumpTarget(flags)
	if len(a) != 6 {
		t.Fatalf("did not select 6 paths (%v)", len(a))
	}
}

func TestBuildDumper(t *testing.T) {
	flags := &appFlags{}

	flags.merge = false
	a := buildDumper(flags)
	if reflect.TypeOf(a).String() != "*wpdump.WPDumper" {
		t.Fatalf("is not WPDump (%v)", reflect.TypeOf(a))
	}

	flags.merge = true
	a = buildDumper(flags)
	if reflect.TypeOf(a).String() != "*wpdump.WPMergeDumper" {
		t.Fatalf("is not MergeDumper (%v)", reflect.TypeOf(a))
	}
}
