package synmedreader

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func processCSV(filename string) ([]sale, error) {
	//Progress Booleans
	var foundFromToDates, foundSalesData bool = false, false
	var storeNum, dinColumn, qtyColumn int = 0, 0, 0
	var currentStoreName, saleDate string = "", ""
	sales := make([]sale, 0)

	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	//This is the row iterator
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			//if there is a random error - just cancel the read
			return sales, err
		}

		var hasConcentration, hasNDC, hasMfr, hasQty, hasCost, hasTotal bool = false, false, false, false, false, false
		var dinFoundForStore bool = false

		din, dinErr := strconv.Atoi(line[dinColumn])
		qty, qtyErr := strconv.ParseFloat(strings.Replace(line[qtyColumn], ",", "", -1), 64)

		// If we are into the sales data already, and we can successfully parse the din and quantity columns, then we're looking pretty good
		if foundFromToDates && foundSalesData && storeNum > 0 && dinErr == nil && qtyErr == nil {
			//The file is formatted in a way where it shows
			//	the total first,
			//	then a breakdown of the transactions
			// --
			// So we're really only interested in the first quantity value for each din for each store

			for iter := 0; iter < len(sales); iter++ {
				if sales[iter].din == din && sales[iter].storeNum == storeNum {
					//Sales data has already been added for this store
					dinFoundForStore = true
				}
			}

			// We didn't find any sales for this store for this din, so let's add it
			if !dinFoundForStore {
				rxNum := makeRxNum(strconv.Itoa(din), storeNum)

				sales = append(sales, sale{
					storeName: currentStoreName,
					saleDate:  saleDate,
					din:       din,
					quantity:  qty,
					storeNum:  storeNum,
					rxNum:     rxNum,
				})
			}

		} else {
			if foundFromToDates {
				// We need to scan for a new store row until the end of the file
				for c := 0; c < int(len(line)); c++ {
					if line[c] == "Concentration" {
						hasConcentration = true
					} else if line[c] == "NDC" || line[c] == "DIN" {
						hasNDC = true
						dinColumn = c
					} else if line[c] == "Mfr" {
						hasMfr = true
					} else if line[c] == "Qty" {
						hasQty = true
						qtyColumn = c
					} else if line[c] == "Cost" {
						hasCost = true
					} else if line[c] == "Total" {
						hasTotal = true
					}
				}

				if hasConcentration && hasNDC && hasMfr && hasQty && hasCost && hasTotal {
					storeNum++
					// Make a Regex to say we only want letters, numbers and spaces
					storeName, err := checkStoreName(line[0])
					if err != nil {
						return sales, err
					}

					currentStoreName = storeName
					//Just in case it's not already true (basically just for the headers before the data begins)
					foundSalesData = true
				}
			} else {
				//Scan row for "To", try to find the to dates
				for c := 0; c < int(len(line)); c++ {
					if line[c] == "To" {
						// now we're going to search the next 5 columns until we find a parsable date - it's probably going to land on the first one, but we're just covering bases
						for i := 1; i < 5; i++ {
							if len(saleDate) < 1 && len(line[int(c+i)]) > 1 {
								//try to parse date
								t, err := dateparse.ParseAny(line[int(c+i)])
								if err == nil {
									// We found a date - let's see confirm with the user that it is the date they want to use
									sellDate, dErr := confirmDate(t.Format("2006-01-02"))
									if dErr != nil {
										return sales, dErr
									}
									saleDate = sellDate
								} else {
									sellDate, dErr := confirmDate(time.Now().Format("2006-01-02"))
									if dErr != nil {
										return sales, dErr
									}
									saleDate = sellDate
								}

								foundFromToDates = true
							}
						}
					}
				}

			}
		}
	}
	return sales, nil
}
