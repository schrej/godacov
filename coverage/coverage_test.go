package coverage

import (
	"fmt"
	"regexp"
	"testing"
)

var testFilename = "github.com/user/repo/a/b/c/1/my_file.go"
var testLine = fmt.Sprintf("%s:16.45,19.24 2 1", testFilename)

func TestRegexpFilename(t *testing.T) {
	r, _ := regexp.Compile(regexpStringFilename)
	result := r.FindStringSubmatch(testFilename)
	if len(result) < 2 {
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
		t.Fatal("invalid mode")
	}
	if ModeSet != result[1] {
		t.Error("expected set mode")
	}
}

func TestCalculatePercentages(t *testing.T) {
	files := map[string]*fileCoverage{
		"a": &fileCoverage{
			numStatements:     12,
			cntStatements:     6,
			coveredStatements: 11,
			lines: map[int]int{
				46: 2,
				55: 2,
				73: 2,
				51: 0,
			},
		},
	}
	total, perFile := calculatePercentages(files)
	if total != 91 {
		t.Error(fmt.Sprintf("expected 91%%, received %v%% \n", total))
	}
	if perFile["a"] != 91 {
		t.Error(fmt.Sprintf("expected 91%%, received %v%% \n", perFile["a"]))
	}
}

func TestParseLine(t *testing.T) {
	regex, _ = compileRegexp()
	report, err := parseLine(testLine)
	if err != nil {
		t.Fatal(err)
	}
	if report.numStatements != 2 {
		t.Error("expected 2 statements")
	}
	if report.cntStatements != 1 {
		t.Error("expected count to be 2")
	}
}

func TestParseLineFailsWithBadFormat(t *testing.T) {
	expectParseLineFails(t, fmt.Sprintf("%s:16.45,19.24 100000000000000000000 1", testFilename))
	expectParseLineFails(t, fmt.Sprintf("%s:16.45,19.24 1 100000000000000000000", testFilename))
	expectParseLineFails(t, fmt.Sprintf("%s:16.sdvsfbvs", testFilename))
}

func expectParseLineFails(t *testing.T, line string) {
	_, err := parseLine(line)
	if err != nil {
		return
	}
	t.Error("expected an error")
}

func TestSkippableLine(t *testing.T) {
	if !isSkippableLine(" ") {
		t.Error("empty lines should be skipped")
	}
	if !isSkippableLine("mode: atomic") {
		t.Error("lines only reporting coverage mode should be skipped")
	}
	if isSkippableLine("something else") {
		t.Error("by default lines should not be skipped")
	}
}
