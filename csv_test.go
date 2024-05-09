package main

import (
	//"strings"
	"reflect"
	"testing"
)

//this should be an integration test and not a unit test i think? as it involves the machine's file system and is dependent on it, hence it is not an independent component of the program
// func TestReadCsv(t *testing.T) {}

func TestShiftCsvSlice(t *testing.T) {
	testCsvSlice := [][]string{
		{"C", "D", "A", "B"},
		{"3", "4", "1", "2"},
	}
	testShiftValue := 2

	expected := [][]string{
		{"A", "B", "C", "D"},
		{"1", "2", "3", "4"},
	}

	result := shiftCsvSlice(testCsvSlice, testShiftValue)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %q, wanted %q", result, expected)
	}
}

func TestCompareCsvSlice(t *testing.T) {
	testSliceOne := [][]string{
		{"A", "B", "C", "D"},
		{"1", "2", "3", "4"},
		{"5", "6", "7", "8"},
	}
	testSliceTwo := [][]string{
		{"A", "B", "C", "D"},
		{"1", "2", "3", "5"},
		{"W", "X", "Y", "Z"},
		{"hello", "this", "is", "test"},
	}

	expected := [][]string{
		{"1", "2", "3", "5"},
		{"W", "X", "Y", "Z"},
		{"hello", "this", "is", "test"},
	}

	result := compareCsvSlice(testSliceOne, testSliceTwo)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %q, wanted %q", result, expected)
	}
}

func TestConvertIntoHash(t *testing.T) {
	testSlice := []string{"this", "will", "be", "joined", "and", "hashed"}

	expected := "1b0105b0d5ce99dfb3d7d0a7b1015b8c"

	result := convertIntoHash(testSlice)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %q, wanted %q", result, expected)
	}
}

//not sure if this is needed to be unit tested as it is a nested inner function solely used by mismatchCheck only
//therefore if mismatchCheck passes all test cases, then appendMismatchRemark has passed as well?
// func TestAppendMismatchRemark(t *testing.T){}

//not sure if this is needed to be unit tested as it is a nested inner function solely used by mismatchCheck only
//therefore if mismatchCheck passes all test cases, then appendMissingRemark has passed as well?
//func TestAppendMissingRemark(t *testing.T){}

func TestMismatchCheck(t *testing.T) {

	testCsvSlice := [][]string{
		{"Amount", "Description", "Date", "ID"},
		{"20", "asd", "2021-01-01", "zxc"},
		{"2", "sss", "2021-01-02", "xcv"},
		{"3", "dfg", "1999-12-01", "cvb"},
		{"4", "fgh", "2021-01-04", "qwerty"},
	}

	testDiffSlice := [][]string{
		{"1", "asd", "2021-01-01", "zxc"},
		{"2", "sdf", "2021-01-02", "xcv"},
		{"3", "dfg", "2021-01-03", "cvb"},
		{"4", "fgh", "2021-01-04", "vbn"},
		{"5", "ghj", "2021-01-05", "bnm"},
	}

	testErrorTolerance := 1

	expected := [][]string{
		{"Amount", "Description", "Date", "ID", "Remarks"},
		{"1", "asd", "2021-01-01", "zxc", "MISMATCH IN COLUMNS: AMOUNT "},
		{"2", "sdf", "2021-01-02", "xcv", "MISMATCH IN COLUMNS: DESCRIPTION "},
		{"3", "dfg", "2021-01-03", "cvb", "MISMATCH IN COLUMNS: DATE "},
		{"4", "fgh", "2021-01-04", "vbn", "MISMATCH IN COLUMNS: ID "},
		{"5", "ghj", "2021-01-05", "bnm", "ERROR: MISSING ENTRY"},
	}

	result := mismatchCheck(testDiffSlice, testCsvSlice, testErrorTolerance)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %q\n", result)
		t.Errorf("Wanted %q\n", expected)
	}

}

//this should be an integration test and not a unit test i think? as it involves the machine's file system and is dependent on it, hence it is not an independent component of the program
//func TestOutputToCsv(t *testing.T){}

func TestReportDateRange(t *testing.T) {
	testInputSlice := [][]string{
		{"Amount", "Description", "Date", "ID", "Remarks"},
		{"1", "asd", "2021-01-01", "zxc", "MISMATCH IN COLUMNS: AMOUNT "},
		{"2", "sdf", "2021-01-02", "xcv", "MISMATCH IN COLUMNS: DESCRIPTION "},
		{"3", "dfg", "2021-01-03", "cvb", "MISMATCH IN COLUMNS: DATE "},
		{"4", "fgh", "2021-01-04", "vbn", "MISMATCH IN COLUMNS: ID "},
		{"5", "ghj", "2021-01-05", "bnm", "ERROR: MISSING ENTRY"},
	}

	expected := []string{
		"2021-01-01",
		"2021-01-05",
	}

	result := reportDateRange(testInputSlice)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %q\n", result)
		t.Errorf("Wanted %q\n", expected)
	}
}
