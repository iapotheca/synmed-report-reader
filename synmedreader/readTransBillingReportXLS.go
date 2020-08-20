package synmedreader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	xlsreader "github.com/extrame/xls"

	"github.com/araddon/dateparse"
)

func processXLS(filename string) ([]sale, error) {

	//Configuration Options
	var attemptCustomDateFirst bool = true
	var customDateFormat string = "02/01/2006"

	//Progress Booleans
	var foundFromToDates, foundSalesData bool = false, false
	var storeNum, dinColumn, qtyColumn int = 0, 0, 0
	var currentStoreName, saleDate string = "", ""

	sales := make([]sale, 0)
	xlsFile, err := xlsreader.Open(filename, "utf-8")
	if err != nil {
		return sales, err
	}
	fmt.Println("Now Processing file: ", filename)

	sheet1 := xlsFile.GetSheet(0)

	if sheet1 == nil {
		return sales, errors.New("File is unreadable - no data in sheet 1")
	}

	//This is the Row Iterator
	for r := 0; r <= (int(sheet1.MaxRow)); r++ {
		row1 := sheet1.Row(r)

		var hasConcentration, hasNDC, hasMfr, hasQty, hasCost, hasTotal bool = false, false, false, false, false, false
		var dinFoundForStore bool = false

		din, dinErr := strconv.Atoi(row1.Col(dinColumn))
		qty, qtyErr := strconv.ParseFloat(strings.Replace(row1.Col(qtyColumn), ",", "", -1), 64)

		// If we are into the sales data already, and we can successfully parse the din and quantity columns, then we're looking pretty good
		if foundFromToDates && foundSalesData && storeNum > 0 && dinErr == nil && qtyErr == nil && din > 0 {

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
				for c := 0; c <= int(row1.LastCol()); c++ {
					if row1.Col(c) == "Concentration" {
						hasConcentration = true
					} else if row1.Col(c) == "NDC" || row1.Col(c) == "DIN" {
						hasNDC = true
						dinColumn = c
					} else if row1.Col(c) == "Mfr" {
						hasMfr = true
					} else if row1.Col(c) == "Qty" {
						hasQty = true
						qtyColumn = c
					} else if row1.Col(c) == "Cost" {
						hasCost = true
					} else if row1.Col(c) == "Total" {
						hasTotal = true
					}
				}

				if hasConcentration && hasNDC && hasMfr && hasQty && hasCost && hasTotal {
					storeNum++
					// Make a Regex to say we only want letters, numbers and spaces
					storeName, err := checkStoreName(row1.Col(0))
					if err != nil {
						return sales, err
					}

					currentStoreName = storeName
					//Just in case it's not already true (basically just for the headers before the data begins)
					foundSalesData = true
				}
			} else {
				//Scan row for "To", try to find the to dates
				for c := 0; c <= int(row1.LastCol()); c++ {
					if row1.Col(c) == "To" {
						// now we're going to search the next 5 columns until we find a parsable date - it's probably going to land on the first one, but we're just covering bases
						for i := 1; i < 5; i++ {

							if len(saleDate) < 1 {
								//Lets first get the value of the row/col
								potentialDate := row1.Col(int(c + i))

								if len(potentialDate) > 6 {
									//Now we want to split the string up into sectionS (if there is a space, use the first section
									testArray := strings.Fields(potentialDate)

									//we want to process the first one of these (if there is more than one

									testableDate := ""

									if len(testArray) > 1 {
										testableDate = testArray[0]
									} else {
										testableDate = potentialDate
									}

									fmt.Println("Now checking out col ", int(c+i), " with value: ", testableDate)

									if attemptCustomDateFirst {
										//Attempt this format first, then attempt parse with dateparse package
										t, err := time.Parse(customDateFormat, testableDate)
										if err == nil {

											sellDate, dErr := confirmDate(t.Format("2006-01-02"))
											if dErr != nil {
												return sales, dErr
											}
											saleDate = sellDate
										}
									}

									if len(saleDate) < 1 {
										//try to parse date
										t, err := dateparse.ParseAny(testableDate)
										if err == nil {
											// We found a date - let's see confirm with the user that it is the date they want to use
											sellDate, dErr := confirmDate(t.Format("2006-01-02"))
											if dErr != nil {
												return sales, dErr
											}
											saleDate = sellDate
										}
									}

									foundFromToDates = true
								}
							}
						}

						if len(saleDate) < 1 {
							sellDate, dErr := confirmDate(time.Now().Format("2006-01-02"))
							if dErr != nil {
								return sales, dErr
							}
							saleDate = sellDate
						}
					}
				}

			}
		}
	}

	return sales, nil
}
