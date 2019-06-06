package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func main() {
	jsonLogs, err := ioutil.ReadFile("1.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonSuites, err := ioutil.ReadFile("2.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonCaptures, err := ioutil.ReadFile("3.json")
	if err != nil {
		fmt.Println(err)
	}

	logs := LogSlice{}
	err = json.Unmarshal(jsonLogs, &logs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Logs: %v\n", logs)

	suites := SuiteSlice{}
	err = json.Unmarshal(jsonSuites, &suites)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Suites: %v\n", suites)

	captures := CaptureSlice{}
	err = json.Unmarshal(jsonCaptures, &captures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Captures: %v\n\n", captures)

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

	Tests := TestSlice{}

	for _, value := range TestMap {
		Tests.Tests = append(Tests.Tests, *value)
	}
	fmt.Printf("Tests structure: %v\n\n", Tests)

	jsonTests, _ := json.Marshal(Tests)
	fmt.Printf("JSON Data: %v\n", string(jsonTests))

	err = ioutil.WriteFile("result.json", jsonTests, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
