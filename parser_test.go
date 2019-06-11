package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"testing"
)

func TestLogValidOK(t *testing.T) {
	if err := Validation("Tests/goodValidLogs.json", "Schemas/logsSchema.json"); err != nil {
		t.Errorf("Logs validation failed. %v", err)
	}
}

func TestSuiteValidOK(t *testing.T) {
	if err := Validation("Tests/goodValidSuites.json", "Schemas/suitesSchema.json"); err != nil {
		t.Errorf("Suites validation failed. %v", err)
	}
	t.Skip()
}

func TestCaptureValidOK(t *testing.T) {
	if err := Validation("Tests/goodValidCaptures.json", "Schemas/capturesSchema.json"); err != nil {
		t.Errorf("Captures validation failed. %v", err)
	}
}

func TestLogValidFail(t *testing.T) {
	if err := Validation("Tests/failValidLogs.json", "Schemas/logsSchema.json"); err == nil {
		t.Errorf("Logs validation failed (wrong JSON). %v", err)
	}
}

func TestLogsAnalizator(t *testing.T) {
	Time := "946684810"
	logs := LogSlice{[]Log{{Time, "Test output A", "fail"}}}
	TestMap, MockMap := make(map[int64]*Test), make(map[int64]*Test)
	LogsAnalizator(logs, TestMap)
	unixTime, _ := strconv.ParseInt(Time, 10, 64)
	MockMap[unixTime] = &Test{Name: "Test output A", Status: "fail"}
	if len(TestMap) != len(MockMap) {
		t.Fatal("LogsAnalizator test failed: different maps len")
	}
	for k, v := range TestMap {
		temp, ok := MockMap[k]
		if !ok || *temp != *v {
			t.Error("LogsAnalizator test failed: different maps values")
		}
	}
}

func TestRun(t *testing.T) {
	cmd := exec.Command("sh", "-c", "go run parser.go Data/1.json Data/2.json Data/3.json Tests/testResult.json")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	goodResult, err := ioutil.ReadFile("Tests/goodResult.json")
	if err != nil {
		log.Fatal(err)
	}
	testResult, err := ioutil.ReadFile("Tests/testResult.json")
	if err != nil {
		log.Fatal(err)
	}
	if len(testResult) != len(goodResult) {
		t.Fatal("Example test failed: different result lens")
	}
	for i := range goodResult {
		if testResult[i] != goodResult[i] {
			t.Fatal("Example test failed: different results")
		}
	}

}
