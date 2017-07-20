package main
import "testing"

func TestMain(t *testing.T) {
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
	for _, pathList := range a {
		if len(pathList) != 2 {
			t.Error("Expected two duplicate files")
		}
	}
}
