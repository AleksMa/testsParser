package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	timeModule "time"
)

type Log struct {
	Time   string `json:"time"`
	Test   string `json:"test"`
	Output string `json:"output"`
}
type LogSlice struct {
	Logs []Log `json:"logs"`
}
type Case struct {
	Name   string `json:"name"`
	Errors int    `json:"errors"`
	Time   string `json:"time"`
}
type Suite struct {
	Name  string `json:"name"`
	Tests int    `json:"tests"`
	Cases []Case `json:"cases"`
}
type SuiteSlice struct {
	Suites []Suite `json:"suites"`
}
type Capture struct {
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Time     string `json:"time"`
}
type CaptureSlice struct {
	Captures []Capture `json:"captures"`
}
type Test struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}
type TestSlice struct {
	Tests []Test `json:"tests"`
}

func outputToStatus(output string) string {
	if output == "fail" {
		return output
	} else {
		return "OK"
	}
}

func errorsToStatus(errors int) string {
	if errors == 0 {
		return "OK"
	} else {
		return "fail"
	}
}

func Read(path string) []byte {
	jsonLogs, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("ERROR: incorrect input file. ", err)
		os.Exit(1)
	}
	return jsonLogs
}

func DecodeJSON(source []byte, dest interface{}) {
	//TODO: валидация
	err := json.Unmarshal(source, dest)
	if err != nil {
		fmt.Println("ERROR: input JSON failed.", err)
		os.Exit(1)
	}
}

func GetData(buffers []interface{}) {
	for i, buf := range buffers {
		jsonLogs := Read(os.Args[i+1])
		DecodeJSON(jsonLogs, buf)
		fmt.Printf("%T: %v\n", buf, buf)
	}
}

func EncodeJSON(testSlice TestSlice) []byte {
	jsonTests, err := json.Marshal(testSlice)
	if err != nil {
		fmt.Println("ERROR: output JSON failed.", err)
		os.Exit(1)
	}
	fmt.Printf("JSON Data: %v\n", string(jsonTests))
	return jsonTests
}

func Write(data []byte) {
	path := "Data/result.json"
	if len(os.Args) > 4 {
		path = os.Args[4]
	}

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		fmt.Println("ERROR: output write failed.", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: \"go run testsParser file1 file2 file3 [result]\"")
		return
	}

	logs, suites, captures := LogSlice{}, SuiteSlice{}, CaptureSlice{}
	buffers := []interface{}{&logs, &suites, &captures}
	GetData(buffers)

	TestMap := make(map[int64]*Test)

	for _, log := range logs.Logs {
		time, err := strconv.ParseInt(log.Time, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		if test, exist := TestMap[time]; exist {
			if test.Name != log.Test || test.Status != outputToStatus(log.Output) {
				fmt.Println("Incorrect data...")
			}
		}
		TestMap[time] = &Test{Name: log.Test, Status: outputToStatus(log.Output)}
	}

	for _, suite := range suites.Suites {
		for _, log := range suite.Cases {
			timeT, _ := timeModule.Parse(timeModule.RFC850, log.Time)
			time := timeT.Unix()
			if test, exist := TestMap[time]; exist {
				if test.Name != log.Name || test.Status != errorsToStatus(log.Errors) {
					fmt.Println("Incorrect data...")
				}
			} else {
				TestMap[time] = &Test{Name: log.Name, Status: errorsToStatus(log.Errors)}
			}
		}
	}

	for _, capture := range captures.Captures {
		timeT, _ := timeModule.Parse(timeModule.RFC3339, capture.Time)
		time := timeT.Unix()
		if test, exist := TestMap[time]; exist {
			test.Expected = capture.Expected
			test.Actual = capture.Actual
		} else {
			TestMap[time] = &Test{Expected: capture.Expected, Actual: capture.Actual}
		}
	}

	testSlice := TestSlice{}

	for _, value := range TestMap {
		testSlice.Tests = append(testSlice.Tests, *value)
	}
	fmt.Printf("Test structure: %v\n\n", testSlice)

	jsonTests := EncodeJSON(testSlice)

	Write(jsonTests)
}
