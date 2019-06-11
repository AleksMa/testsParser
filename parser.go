package main

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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

func Validation(path string, schemaPath string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://Schemas/" + schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return err
	}
}

func Read(path string) ([]byte, error) {
	jsonLogs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonLogs, nil
}

func DecodeJSON(source []byte, dest interface{}) error {
	err := json.Unmarshal(source, dest)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJSON(testSlice TestSlice) ([]byte, error) {
	jsonTests, err := json.Marshal(testSlice)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JSON Data: %v\n", string(jsonTests))
	return jsonTests, nil
}

func Write(data []byte) error {
	path := "Data/result.json"
	if len(os.Args) > 4 {
		path = os.Args[4]
	}

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func logsAnalizator(logs LogSlice, TestMap map[int64]*Test) {
	for _, tempLog := range logs.Logs {
		unixTime, err := strconv.ParseInt(tempLog.Time, 10, 64)
		if err != nil {
			log.Fatal("ERROR: Log time parse failed. ", err)
		}
		if test, exist := TestMap[unixTime]; exist {
			if test.Name != tempLog.Test || test.Status != outputToStatus(tempLog.Output) {
				fmt.Println("WARNING: Two different logs with same time")
			}
		}
		TestMap[unixTime] = &Test{Name: tempLog.Test, Status: outputToStatus(tempLog.Output)}
	}
}

func suiteAnalizator(suites SuiteSlice, TestMap map[int64]*Test) {
	for _, suite := range suites.Suites {
		for _, tempLog := range suite.Cases {
			tempTime, err := time.Parse(time.RFC850, tempLog.Time)
			if err != nil {
				log.Fatal("ERROR: Suite time parse failed. ", err)
			}
			unixTime := tempTime.Unix()
			if test, exist := TestMap[unixTime]; exist {
				if test.Name != tempLog.Name || test.Status != errorsToStatus(tempLog.Errors) {
					fmt.Println("WARNING: Two different logs with same time")
				}
			}
			TestMap[unixTime] = &Test{Name: tempLog.Name, Status: errorsToStatus(tempLog.Errors)}
		}
	}
}

func capturesAnalizator(captures CaptureSlice, TestMap map[int64]*Test) {
	for _, capture := range captures.Captures {
		tempTime, err := time.Parse(time.RFC3339, capture.Time)
		if err != nil {
			log.Fatal("ERROR: Capture time parse failed. ", err)
		}
		unixTime := tempTime.Unix()
		if test, exist := TestMap[unixTime]; exist {
			test.Expected = capture.Expected
			test.Actual = capture.Actual
		} else {
			fmt.Println("WARNING: Test with no name and status")
			TestMap[unixTime] = &Test{Expected: capture.Expected, Actual: capture.Actual}
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: \"go run parser.go file1 file2 file3 [fileResult]\"")
		return
	}

	logs, suites, captures := LogSlice{}, SuiteSlice{}, CaptureSlice{}
	buffers := []interface{}{&logs, &suites, &captures}

	schemas := []string{"logsSchema.json", "suitesSchema.json", "capturesSchema.json"}

	for i, buf := range buffers {
		err := Validation(os.Args[i+1], schemas[i])
		if err != nil {
			log.Fatal("ERROR: json validation failed. ", err)
		}
		jsonLogs, err := Read(os.Args[i+1])
		if err != nil {
			log.Fatal("ERROR: input read failed. ", err)
		}
		err = DecodeJSON(jsonLogs, buf)
		if err != nil {
			log.Fatal("ERROR: input JSON failed. ", err)
		}
		fmt.Printf("%T: %v\n", buf, buf)
	}

	TestMap := make(map[int64]*Test)

	logsAnalizator(logs, TestMap)
	suiteAnalizator(suites, TestMap)
	capturesAnalizator(captures, TestMap)

	testSlice := TestSlice{}
	for _, value := range TestMap {
		testSlice.Tests = append(testSlice.Tests, *value)
	}
	fmt.Printf("Test structure: %v\n\n", testSlice)

	jsonTests, err := EncodeJSON(testSlice)
	if err != nil {
		log.Fatal("ERROR: output JSON failed.", err)
	}

	err = Write(jsonTests)
	if err != nil {
		log.Fatal("ERROR: output write failed.", err)
	}
}
