package coverage

import (
	"fmt"
	"regexp"
	"testing"
)

var testFilename = "github.com/user/repo/a/b/c/1/my_file.go"

func TestRegexpFilename(t *testing.T) {
	r, _ := regexp.Compile(regexpStringFilename)
	result := r.FindStringSubmatch(testFilename)
	if len(result) < 2 {
		fmt.Println(result)
		t.Fatal("filename should match")
	}
	if testFilename != result[1] {
		t.Error("filename is not valid")
	}
}

func TestRegexpMode(t *testing.T) {
	modeString := "mode: set"
	r, _ := regexp.Compile(regexpStringMode)
	result := r.FindStringSubmatch(modeString)
	if len(result) < 2 {
		fmt.Println(result)
		t.Fatal("invalid mode")
	}
	if ModeSet != result[1] {
		t.Error("expected set mode")
	}
}
