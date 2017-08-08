package coverage

import (
	"fmt"
	"regexp"
	"testing"
)

var testFilename = "github.com/user/repo/a/b/c/1/my_file.go"

func TestRegexp(t *testing.T) {
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
