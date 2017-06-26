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
	numStatements int
	cntStatements int
	lines         map[int]int
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

var regex *regexp.Regexp

// GenerateCoverageJSON generates a json string containing
// coverage information in codacy's format
func GenerateCoverageJSON(coverageFile string) ([]byte, error) {
	regex, _ = regexp.Compile(`([a-zA-Z\/\.]*):(\d*)\..* (\d*) (\d*)`)

	dat, err := ioutil.ReadFile(coverageFile)
	lines := strings.Split(string(dat), "\n")
	if err != nil {
		return nil, err
	}

	files := make(map[string]*fileCoverage)
	for _, line := range lines[1 : len(lines)-1] {
		parsed, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		//fmt.Println(parsed)
		file := files[parsed.file]
		if file == nil {
			file = new(fileCoverage)
			files[parsed.file] = file
			file.lines = make(map[int]int)
		}
		file.cntStatements += parsed.cntStatements
		file.numStatements += parsed.numStatements
		file.lines[parsed.line] += parsed.cntStatements
	}

	total, perFile := calculatePercentages(files)
	//fmt.Println(total, perFile)

	covJSON := codacyCoverageJSON{}
	covJSON.Total = int(total * 100)
	covJSON.FileReports = make([]codacyFileCoverageJSON, 0)

	for filename, fileCoverage := range perFile {
		fileCov := codacyFileCoverageJSON{}
		fileCov.Filename = filename
		fileCov.Total = int(fileCoverage * 100)
		fileCov.Coverage = files[filename].lines
		covJSON.FileReports = append(covJSON.FileReports, fileCov)
	}

	json, err := json.Marshal(covJSON)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(json))

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

func calculatePercentages(files map[string]*fileCoverage) (float64, map[string]float64) {
	totalNumStatements := 0
	totalCntStatements := 0
	percentages := make(map[string]float64)

	for file, coverage := range files {
		totalNumStatements += coverage.numStatements
		totalCntStatements += coverage.cntStatements
		percentages[file] = calculatePercentage(coverage.numStatements, coverage.cntStatements)
	}

	return calculatePercentage(totalNumStatements, totalCntStatements), percentages
}

func calculatePercentage(num int, cnt int) float64 {
	if num == 0 {
		return 0
	}
	return float64(cnt) / float64(num)
}
