//APPLICANT: Saiful Amin

package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	//"reflect"
)

// function to read csv
func readCsv(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	lines, err := reader.ReadAll()
	if err == io.EOF {
		fmt.Println(err)
	}

	return lines
}

// function to rotate either one of the CSVs by N positions to the left
// this will ensure that both csv slices will have the same columns instead of mismatched columns
// accepts one slice of string slices as argument 1 and an integer value shiftValue to shift to the left
func shiftCsvSlice(csvSlice [][]string, shiftValue int) [][]string {

	for j := 0; j < shiftValue; j++ {
		for i := 0; i < len(csvSlice); i++ {
			var x string
			x, csvSlice[i] = csvSlice[i][0], csvSlice[i][1:] //pops off the first element of inner slice
			csvSlice[i] = append(csvSlice[i], x)             //appends the previously popped off element to inner slice
		}
	}

	return csvSlice
}

// function to compare the two slices and output the difference as a slice of string slices
func compareCsvSlice(sliceOne, sliceTwo [][]string) [][]string {

	var diff [][]string
	sliceHash := make(map[string][]string)
	//slices cannot be used as key for map, hence we will get hex string by hashing and encoding to use as key

	//need to check lengths of slices and ensure the longer slice is used when checking against existence of value in hash table
	//if we use the shorter slice for difference appending we will miss some values
	if len(sliceOne) > len(sliceTwo) {

		for _, element := range sliceTwo[1:] {
			sliceHash[convertIntoHash(element)] = element
		}

		for _, element := range sliceOne[1:] {
			if _, exists := sliceHash[convertIntoHash(element)]; !exists {
				diff = append(diff, element)
			}
		}

	} else if len(sliceOne) < len(sliceTwo) {

		for _, element := range sliceOne[1:] {
			sliceHash[convertIntoHash(element)] = element
		}

		for _, element := range sliceTwo[1:] {
			if _, exists := sliceHash[convertIntoHash(element)]; !exists {
				diff = append(diff, element)
			}
		}

	} else {

		for _, element := range sliceOne[1:] {
			sliceHash[convertIntoHash(element)] = element
		}

		for _, element := range sliceTwo[1:] {
			if _, exists := sliceHash[convertIntoHash(element)]; !exists {
				diff = append(diff, element)
			}
		}
	}

	return diff
}

// function to join string slice into a single string and then md5 hash into byte array before
// encoding it in a base-16 string for use in compareCsvSlice function for key in hashmap
func convertIntoHash(slice []string) string {

	var joinedString string
	var hashByteArray [16]byte
	var hexEncodedString string
	joinedString = strings.Join(slice, "")
	//fmt.Println("JOINED STRING: ", joinedString)
	hashByteArray = md5.Sum([]byte(joinedString))
	//fmt.Println("HASH BYTE ARRAY: ", hashByteArray)
	hexEncodedString = hex.EncodeToString(hashByteArray[:])
	//fmt.Println("HEX ENCODED STRING: ", hexEncodedString)
	return hexEncodedString
}

// function to append the mismatch remarks such as mismatched amount, descrption, id, etc to the remarks column in csv
// we are updating mismatchedOutput by passing it as a pointer
// function block moved out of mismatchRemark
func appendMismatchRemark(csvSlice [][]string, diffSlice [][]string, mismatchedOutput *[][]string, csvSliceElem int, diffSliceElem int, commonElem int) {

	switch {
	case ((csvSlice[csvSliceElem][commonElem] != diffSlice[diffSliceElem][commonElem]) && commonElem == 0):
		diffSlice[diffSliceElem][len(diffSlice[diffSliceElem])-1] += "AMOUNT "
		*mismatchedOutput = append(*mismatchedOutput, diffSlice[diffSliceElem])
	case ((csvSlice[csvSliceElem][commonElem] != diffSlice[diffSliceElem][commonElem]) && commonElem == 1):
		diffSlice[diffSliceElem][len(diffSlice[diffSliceElem])-1] += "DESCRIPTION "
		*mismatchedOutput = append(*mismatchedOutput, diffSlice[diffSliceElem])
	case ((csvSlice[csvSliceElem][commonElem] != diffSlice[diffSliceElem][commonElem]) && commonElem == 2):
		diffSlice[diffSliceElem][len(diffSlice[diffSliceElem])-1] += "DATE "
		*mismatchedOutput = append(*mismatchedOutput, diffSlice[diffSliceElem])
	case ((csvSlice[csvSliceElem][commonElem] != diffSlice[diffSliceElem][commonElem]) && commonElem == 3):
		diffSlice[diffSliceElem][len(diffSlice[diffSliceElem])-1] += "ID "
		*mismatchedOutput = append(*mismatchedOutput, diffSlice[diffSliceElem])
	}
}

// function to append the missing entries with a remark noting it as such
// we are updating mismatchedOutput by passing it as a pointer
// function block moved out of mismatchRemark
func appendMissingRemark(diffSlice [][]string, mismatchedOutput *[][]string, diffSliceElem int) {

	diffSlice[diffSliceElem] = append(diffSlice[diffSliceElem], "ERROR: MISSING ENTRY")
	*mismatchedOutput = append(*mismatchedOutput, diffSlice[diffSliceElem])
}

// function used to check which columns of the diffSlice and csvSlice are not matching and output the rows.
// requires error tolerance integer as we are comparing element by element of the two arrays, and then column by column of each inner array, hence some elements of diffSlice and elements of csvSlice are totally different and will provide a false(?) negative
// with errorTolerance argument, we can specify how many errors we want to spot and in which columns.
// example: with errorTolerance = 1, we can spot which elements in the two argument slices have only one mismatch. This is the output we want in order to accurately see mismatched data.
// with errorTolerance = 5 however, we will get a lot of false negative such as comparing diffSlice[0] which is [30 D 2021-07-03 zodo] and csvSlice[1] which is [24 A 2021-06-30 zoUr].
// BOTH are 5/5 errors, BUT [24 A 2021-06-30 zoUr] exists in proxy.csv, so it is a false negative
// function is left as open ended as possible to allow user to input any error tolerance they want instead of me hardcoding methods such as checking against user ID etc, as I am going with the assumption that we do not
// know WHAT column is correct and what is not
// this of course comes with the disadvantage that this function is O(n^4) at worse case :(
func mismatchCheck(diffSlice [][]string, csvSlice [][]string, errorTolerance int) [][]string {

	errorCounter := 0
	var mismatchedOutput [][]string
	mismatchedOutput = append(mismatchedOutput, csvSlice[0])     //get headers
	mismatchedOutput[0] = append(mismatchedOutput[0], "Remarks") //add Remarks column to mismatchedOutput

	for i := 0; i < len(csvSlice); i++ {
		for j := 0; j < len(diffSlice); j++ {
			for k := 0; k < len(diffSlice[j]); k++ {
				if diffSlice[j][k] != csvSlice[i][k] {
					errorCounter++
				}
			}

			if errorCounter == errorTolerance {

				diffSlice[j] = append(diffSlice[j], "MISMATCH IN COLUMNS: ")

				for a := 0; a < len(csvSlice[i]); a++ {

					appendMismatchRemark(csvSlice, diffSlice, &mismatchedOutput, i, j, a)

				}

				copy(diffSlice[j:], diffSlice[j+1:])
				diffSlice[len(diffSlice)-1] = nil
				diffSlice = diffSlice[:len(diffSlice)-1]

			} else {
				errorCounter = 0
			}
		}
	}

	for i := 0; i < len(diffSlice); i++ {
		appendMissingRemark(diffSlice, &mismatchedOutput, i)
	}

	return mismatchedOutput
}

func outputToCsv(inputSlice [][]string) {
	outputCsv, err := os.Create("output.csv")

	if err != nil {
		log.Fatalln("Error creating CSV file: ", err)
	}
	wr := csv.NewWriter(outputCsv)

	for _, row := range inputSlice {
		_ = wr.Write(row)
	}

	//wr.WriteAll(inputSlice)
	if err := wr.Error(); err != nil {
		log.Fatalln("Error writing CSV: ", err)
	}
	wr.Flush()
	outputCsv.Close()
}

func reportDateRange(inputSlice [][]string) []string {
	tempSlice := inputSlice[1:]

	sort.Slice(tempSlice, func(i, j int) bool {
		return tempSlice[i][2] < tempSlice[j][2]
	})

	var dateRange []string

	dateRange = append(dateRange, tempSlice[0][2])
	dateRange = append(dateRange, tempSlice[len(tempSlice)-1][2])
	// layout := "2006-01-02"
	// earliestDate, err := time.Parse(layout, inputSlice[1][2])
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// latestDate, err := time.Parse(layout, inputSlice[len(inputSlice)-1][2])
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// dateRange := latestDate.Sub(earliestDate)
	// fmt.Println("Time between: ", dateRange)

	return dateRange
}
