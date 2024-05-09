//APPLICANT: Saiful Amin

package main

func main() {
	proxyCsv := readCsv("proxy.csv")   //length 11
	sourceCsv := readCsv("source.csv") //length 10
	shiftCsvSlice(sourceCsv, 2)        //first argument: target csv to shift to the left by N integers, second argument: value of N

	diff := compareCsvSlice(proxyCsv, sourceCsv)
	mismatchedData := mismatchCheck(diff, sourceCsv, 1)
	outputToCsv(mismatchedData)

	//fmt.Println(reportDateRange(mismatchedData))
	outputTextReport(reportDateRange(mismatchedData), sourceCsv, mismatchedData)
}
