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

func rm(p string) error {
	fmt.Printf("pretending to delete during the test '%s'\n", p)
	return nil
}

func TestPurge(t *testing.T) {
	s := trait{2, []string{"test/c", "test/a/c"}}
	if s.purge(false, "", rm) != 0 {
		t.Error("something deleted when prefix option is empty with verbose flag off")
	}

	if s.purge(true, "", rm) != 0 {
		t.Error("something deleted when prefix option is empty with verbose flag true")
	}

	if s.purge(true, "test/a", rm) != 1 {
		t.Error("nothing is deleted when one file is expected to be gone")
	}
}
