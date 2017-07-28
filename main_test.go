package main

import (
	"fmt"
	"testing"
)

func TestHardID(t *testing.T) {
	a := hardID([]string{"test/e", "test/d"})
	if len(a) != 2 {
		t.Error("Expected no dupes")
	}
	for _, pathList := range a {
		if len(pathList) != 1 {
			t.Error("Expected only one file in the list")
		}
	}

	a = hardID([]string{"test/c", "test/d"})
	if len(a) != 1 {
		t.Error("Expected two files to have the same fingerprint")
	}
	for m, pathList := range a {
		if len(pathList) != 2 {
			t.Error("Expected two duplicate files")
		}
		if m != "b026324c6904b2a9cb4b88d6d61c81d1" {
			t.Error("Expected md5 sum must not be %v", m)
		}
	}
	a = hardID([]string{"test/a/b", "test/a/d"})
	if len(a) != 1 {
		t.Error("Expected only one valid file")
	}
}

func TestConfirmDupes(t *testing.T) {
	s := trait{0, []string{}}
	if s.confirmDupes(true) {
		t.Error("Dupes found when none expected for trivial trait")
	}

	s = trait{0, []string{"a", "b"}}
	if s.confirmDupes(true) {
		t.Error("Dupes found for 0 size files")
	}

	s = trait{9, []string{"/root"}}
	if s.confirmDupes(true) {
		t.Error("Dupes report when there is only one copy")
	}

	s = trait{2, []string{"test/c", "test/a/c"}}
	if !s.confirmDupes(true) {
		t.Error("Dupes are not reported")
	}

	s = trait{3, []string{"test/e", "test/a/e"}}
	if s.confirmDupes(false) {
		t.Error("Dupes are unexpectedly reported mith quiet setting off")
	}

	if s.confirmDupes(true) {
		t.Error("Dupes are unexpectedly reported with quiet setting")
	}
}

func TestCompileData(t *testing.T) {
	s, c, data := compileData([]string{"test"})
	if s != 15 {
		t.Errorf("total size is expected to be %d but got %d", 15, s)
	}

	if c != 6 {
		t.Errorf("count of files is expected to be %d but instead %d", 6, c)
	}

	if len(data) != 4 {
		t.Errorf("expected data length is %d but actual value is %d", 4, len(data))
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

	if s.purge(true, "test/a", rm) == 0 {
		t.Error("nothing is deleted when one file is expected to be gone")
	}
}
