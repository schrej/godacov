package coverage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type reportLine struct {
	file          string
	line          int
	numStatements int
	cntStatements int
}

type fileCoverage struct {
	numStatements     int
	cntStatements     int
	coveredStatements int
	lines             map[int]int
}

type codacyCoverageJSON struct {
	Total       int                      `json:"total"`
	FileReports []codacyFileCoverageJSON `json:"fileReports"`
}

type codacyFileCoverageJSON struct {
	Filename string      `json:"filename"`
	Total    int         `json:"total"`
	Coverage map[int]int `json:"coverage"`
}

const (
	ModeSet = "set"
)

var regex *regexp.Regexp
var regexpStringFilename = `([a-zA-Z\/\._\d]*)`
var regexpStringStat = `(\d*)`
var regexpStringMode = `mode: ([set|count|atomic]*)`
var regexpString = fmt.Sprintf(`%s:%s\..* %s %s`, regexpStringFilename, regexpStringStat, regexpStringStat, regexpStringStat)

// GenerateCoverageJSON generates a json string containing
// coverage information in codacy's format
func GenerateCoverageJSON(coverageFile string) ([]byte, error) {
	regex, _ = regexp.Compile(regexpString)

	dat, err := ioutil.ReadFile(coverageFile)
	lines := strings.Split(string(dat), "\n")
	if err != nil {
		return nil, err
	}

	files := make(map[string]*fileCoverage)
	for _, line := range lines {
		if isSkippableLine(line) {
			continue
		}
		parsed, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		file := files[parsed.file]
		if file == nil {
			file = new(fileCoverage)
			files[parsed.file] = file
			file.lines = make(map[int]int)
		}
		file.cntStatements += parsed.cntStatements
		file.numStatements += parsed.numStatements
		if parsed.cntStatements > 0 {
			file.coveredStatements += parsed.numStatements
		}
		file.lines[parsed.line] += parsed.cntStatements
	}

	total, perFile := calculatePercentages(files)

	covJSON := codacyCoverageJSON{}
	covJSON.Total = total
	covJSON.FileReports = make([]codacyFileCoverageJSON, 0)

	for filename, fileCoverage := range perFile {
		fileCov := codacyFileCoverageJSON{}
		fileCov.Filename = filename
		fileCov.Total = fileCoverage
		fileCov.Coverage = files[filename].lines
		covJSON.FileReports = append(covJSON.FileReports, fileCov)
	}

	json, err := json.Marshal(covJSON)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func parseLine(line string) (reportLine, error) {
	result := regex.FindStringSubmatch(line)
	if len(result) >= 5 {
		line, err := strconv.Atoi(result[2])
		if err != nil {
			return reportLine{}, err
		}
		numStatements, err := strconv.Atoi(result[3])
		if err != nil {
			return reportLine{}, err
		}
		cntStatements, err := strconv.Atoi(result[4])
		if err != nil {
			return reportLine{}, err
		}

		return reportLine{result[1], line, numStatements, cntStatements}, nil
	}

	return reportLine{}, errors.New("Invalid line format")
}

func isSkippableLine(line string) bool {
	return (len(strings.TrimSpace(line)) == 0) || strings.HasPrefix(line, "mode")
}

func calculatePercentages(files map[string]*fileCoverage) (int, map[string]int) {
	totalNumStatements := 0
	totalCntStatements := 0
	totalCoveredStatements := 0
	percentages := make(map[string]int)

	for file, coverage := range files {
		totalNumStatements += coverage.numStatements
		totalCntStatements += coverage.cntStatements
		totalCoveredStatements += coverage.coveredStatements
		percentages[file] = calculatePercentage(coverage.numStatements, coverage.coveredStatements)
	}

	return calculatePercentage(totalNumStatements, totalCoveredStatements), percentages
}

func calculatePercentage(num int, cvd int) int {
	if num == 0 {
		return 0
	}

	return cvd * 100 / num
}

func compileRegexp() (*regexp.Regexp, error) {
	return regexp.Compile(regexpString)
}
