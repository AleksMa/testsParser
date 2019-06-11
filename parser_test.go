package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	cmd := exec.Command("sh", "-c", "go run parser.go Data/1.json Data/2.json Data/3.json Tests/testResult.json")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	goodResult, err := ioutil.ReadFile("Tests/goodResult1.json")
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
