//APPLICANT: Saiful Amin

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func outputTextReport(dateRange []string, numRecords [][]string, mismatchedData [][]string) {

	f, err := os.Create("textreport.txt")
	if err != nil {
		log.Fatalln("Error creating text file: ", err)
	}

	defer f.Close()

	var dateRangeString = fmt.Sprintf("Date Range for Report is from %v to %v\n", dateRange[0], dateRange[1])
	var sourceRecordsProcessed = fmt.Sprintf("Number of Source Records Processed is: %v\n", len(numRecords))
	var numOfDiscrepancies = fmt.Sprintf("Number of Discrepancies is: %v\n", len(mismatchedData)-1)

	_, err = f.WriteString(dateRangeString)
	_, err = f.WriteString(sourceRecordsProcessed)
	_, err = f.WriteString(numOfDiscrepancies)
	_, err = f.WriteString("The discrepancies are as follows: \n")

	discrepanciesMap := make(map[string]int)
	discrepanciesMap["MISMATCHED AMOUNT"] = 0
	discrepanciesMap["MISMATCHED DESCRIPTION"] = 0
	discrepanciesMap["MISMATCHED DATE"] = 0
	discrepanciesMap["MISMATCHED ID"] = 0
	discrepanciesMap["MISSING ENTRIES"] = 0

	for i := 0; i < len(mismatchedData); i++ {

		discrepanciesString := mismatchedData[i][4]

		if strings.Contains(discrepanciesString, "AMOUNT") {
			discrepanciesMap["MISMATCHED AMOUNT"]++
		}
		if strings.Contains(discrepanciesString, "DESCRIPTION") {
			discrepanciesMap["MISMATCHED DESCRIPTION"]++
		}
		if strings.Contains(discrepanciesString, "DATE") {
			discrepanciesMap["MISMATCHED DATE"]++
		}
		if strings.Contains(discrepanciesString, "ID") {
			discrepanciesMap["MISMATCHED ID"]++
		}
		if strings.Contains(discrepanciesString, "MISSING ENTRY") {
			discrepanciesMap["MISSING ENTRIES"]++
		}
	}

	_, err = f.WriteString(fmt.Sprintf("Number of MISMATCHED AMOUNT: %v\n", discrepanciesMap["MISMATCHED AMOUNT"]))
	_, err = f.WriteString(fmt.Sprintf("Number of MISMATCHED DESCRIPTION: %v\n", discrepanciesMap["MISMATCHED DESCRIPTION"]))
	_, err = f.WriteString(fmt.Sprintf("Number of MISMATCHED DATE: %v\n", discrepanciesMap["MISMATCHED DATE"]))
	_, err = f.WriteString(fmt.Sprintf("Number of MISMATCHED ID: %v\n", discrepanciesMap["MISMATCHED ID"]))
	_, err = f.WriteString(fmt.Sprintf("Number of MISSING ENTRIES: %v\n", discrepanciesMap["MISSING ENTRIES"]))
}
