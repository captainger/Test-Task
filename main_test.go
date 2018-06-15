package main

import (
	"bytes"
	"strings"
	"testing"
)

type testCase struct {
	origins  string
	response []string
	source   string
}

var TestCases = []testCase{
	testCase{
		origins: "TestData\\Test1.txt\nTestData\\Test2.txt\nTestData\\Test3.txt",
		response: []string{
			"Count for TestData\\Test3.txt: 0",
			"Count for TestData\\Test1.txt: 3",
			"Count for TestData\\Test2.txt: 7",
			"Total: 10",
		},
		source: "file",
	},
	testCase{
		origins: "https://golang.org\nhttps://golang.org\nhttps://en.wikipedia.org/wiki/Go_Go_Gophers",
		response: []string{
			"Count for https://en.wikipedia.org/wiki/Go_Go_Gophers: 174",
			"Count for https://golang.org: 9",
			"Count for https://golang.org: 9",
			"Total: 192",
		},
		source: "url",
	},
	testCase{
		origins: "TestData\\Test1.txt\nTestData\\Test2.txt\nTestData\\Test3.txt\nTestData\\Test1.txt\nTestData\\Test2.txt\nTestData\\Test3.txt",
		response: []string{
			"Count for TestData\\Test3.txt: 0",
			"Count for TestData\\Test1.txt: 3",
			"Count for TestData\\Test1.txt: 3",
			"Count for TestData\\Test2.txt: 7",
			"Count for TestData\\Test2.txt: 7",
			"Count for TestData\\Test3.txt: 0",
			"Total: 20",
		},
		source: "file",
	},
	testCase{
		origins: "https://golang.org\nhttps://golang.org\nhttps://en.wikipedia.org/wiki/Go_Go_Gophers\nhttps://www.google.ru/\nhttps://golang.org/pkg/\nhttps://golang.org/help/",
		response: []string{
			"Count for https://www.google.ru/: 10",
			"Count for https://en.wikipedia.org/wiki/Go_Go_Gophers: 174",
			"Count for https://golang.org/pkg/: 33",
			"Count for https://golang.org: 9",
			"Count for https://golang.org/help/: 36",
			"Count for https://golang.org: 9",
			"Total: 271"},
		source: "url",
	},
	testCase{
		origins: "TestData\\Test.txt",
		response: []string{
			"File  TestData\\Test.txt  can not be open",
			"Total: 0",
		},
		source: "file",
	},
	testCase{
		origins: "https://golang.o",
		response: []string{
			"URL  https://golang.o  can not be open",
			"Total: 0"},
		source: "url",
	},
}

func TestGoSearcher(t *testing.T) {
	for i, item := range TestCases {
		input := new(bytes.Buffer)
		output := new(bytes.Buffer)
		input.WriteString(item.origins)
		GoSearcher(input, output, item.source)
		result := strings.Split(output.String(), "\n")
		if len(result) != len(item.response) {
			t.Errorf("Test[%d] is failed: Result has invalid length\n Got:%d\nExpected:%d\n", i, len(result), len(item.response))
		}
		for _, expected := range item.response {
			isEqual := false
			for _, item := range result {
				if item == expected {
					isEqual = true
				}
			}
			if !isEqual {
				t.Errorf("Test[%d] is failed: Got:%v\nExpected:%v\n", i, result, item.response)
			}
		}
	}
}
