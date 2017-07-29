package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHardID(t *testing.T) {
	cases := []struct {
		paths     []string
		errString string
		md5s      map[string][]string
	}{
		{[]string{"test/e", "test/d"},
			"Case: no dupes with two valid files: expected %v got %v",
			map[string][]string{
				"2737b49252e2a4c0fe4c342e92b13285": []string{"test/e"},
				"b026324c6904b2a9cb4b88d6d61c81d1": []string{"test/d"},
			},
		},
		{[]string{"test/c", "test/d"},
			"Case: two good files with same content: expected %v got %v",
			map[string][]string{
				"b026324c6904b2a9cb4b88d6d61c81d1": []string{"test/c", "test/d"},
			},
		},
		{[]string{"test/a/b", "test/a/f"},
			"Case: only one valid file: expected %v got %v",
			map[string][]string{
				"7e2fe280d0a014cf5035bd8dddf59410": []string{"test/a/f"},
			},
		},
	}
	for _, s := range cases {
		h := hardID(s.paths)
		if !reflect.DeepEqual(s.md5s, h) {
			t.Errorf(s.errString, s.md5s, h)
		}
	}
}

func TestConfirmDupes(t *testing.T) {
	cases := []struct {
		t       trait
		errCase bool
		err     string
	}{
		{trait{0, []string{}},
			true,
			"Dupes found when none expected with empty list"},
		{trait{0, []string{"a", "b"}},
			true,
			"Dupes found for zero size files"},
		{trait{9, []string{"test/c"}},
			true,
			"Dupes reported when there is only one copy"},
		{trait{2, []string{"test/c", "test/a/c"}},
			false,
			"Dupers are not reported"},
		{trait{3, []string{"test/e", "test/a/e"}},
			true,
			"Dupes are unexpectedly reported"},
	}

	for _, c := range cases {
		for _, q := range [2]bool{true, false} {
			if c.errCase == c.t.confirmDupes(q) {
				t.Errorf(c.err+" with quiet setting %t\n", q)
			}
		}
	}
}

func TestDoWork(t *testing.T) {
	actualCount, actualSize, actualDupes, actualSaved := doWork(true, true, "test/a", []string{"test"})
	count, size, dupes, saved := 7, int64(18), 1, int64(2)
	if dupes != actualDupes {
		t.Errorf("total dupes expected %d got %d", dupes, actualDupes)
	}

	if count != actualCount {
		t.Errorf("total count expected %d got %d", count, actualCount)
	}

	if size != actualSize {
		t.Errorf("size expected %d got %d", size, actualSize)
	}
	if saved != actualSaved {
		t.Errorf("saved expected %d bytes got %d", saved, actualSaved)
	}
}

func TestCompileData(t *testing.T) {
	s, c, data := compileData([]string{"test"})
	total, count, dataLength := int64(18), 7, 5
	if s != total {
		t.Errorf("total size is expected to be %d but got %d", total, s)
	}

	if c != count {
		t.Errorf("count of files is expected to be %d but instead %d", count, c)
	}

	if len(data) != dataLength {
		t.Errorf("expected data length is %d but actual value is %d", dataLength, len(data))
	}
}

func TestValidateDirs(t *testing.T) {
	if validateDirs([]string{"nodir"}) {
		t.Error("validate non-existing dir")
	}

	if !validateDirs([]string{".", "test"}) {
		t.Error("doest not validate good dirs")
	}
}

type fakeDeletePath struct {
	count int
}

func (f *fakeDeletePath) rm(p string) error {
	fmt.Printf("pretending to delete during the test '%s'\n", p)
	f.count++
	return nil
}

func TestPurge(t *testing.T) {
	s := trait{2, []string{"test/c", "test/a/c"}}
	f := fakeDeletePath{0}
	type testResult struct {
		deleted       int
		actualDeleted int
		err           string
	}
	cases := []struct {
		prefix  string
		results []testResult
	}{
		{"", []testResult{
			{0, 0, "something deleted when prefix option is empty"},
			{0, 0, "something deleted when prefix option is empty"}}},
		{"test/a", []testResult{
			{1, 0, "nothing must be deleted for valid prefix"},
			{1, 1, "exactly one file is expected to be deleted"}}}}
	for _, c := range cases {
		for i, q := range [2]bool{true, false} {
			deleted := s.purge(q, c.prefix, f.rm)
			if deleted != c.results[i].deleted {
				t.Errorf(c.results[i].err+" with report flag %t"+": deleted %d", q, deleted)
			}
			if c.results[i].actualDeleted != f.count {
				t.Errorf("fake rm reports different count %d from code: %d", f.count, deleted)
			}
			f.count = 0
		}
	}
}
