package srs

import (
	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(before, after string) []dmp.Diff {
	d := dmp.New()
	diffs := d.DiffMain(before, after, true)
	return d.DiffCleanupSemantic(diffs)
}
