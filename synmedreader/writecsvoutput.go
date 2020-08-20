package synmedreader

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tushar2708/altcsv"
)

func MakeCSV(sales []sale) error {

	csvfile, err := os.Create("CentralSales.csv")

	if err != nil {
		return err
	}

	csvwriter := altcsv.NewWriter(csvfile)
	csvwriter.AllQuotes = true

	headerRow := make([]string, 0)
	headerRow = append(headerRow, "din")
	headerRow = append(headerRow, "filldate")
	headerRow = append(headerRow, "qty")
	headerRow = append(headerRow, "rx")
	headerRow = append(headerRow, "patientname")

	err = csvwriter.Write(headerRow)
	if err != nil {
		return err
	}

	for s := 0; s < len(sales); s++ {
		dataRow := make([]string, 0)
		dataRow = append(dataRow, strconv.Itoa(sales[s].din))
		dataRow = append(dataRow, sales[s].saleDate)
		dataRow = append(dataRow, strconv.FormatFloat(sales[s].quantity, 'f', 1, 64))
		dataRow = append(dataRow, sales[s].rxNum)
		dataRow = append(dataRow, sales[s].storeName)
		err = csvwriter.Write(dataRow)
		if err != nil {
			return err
		}
	}
	csvwriter.Flush()

	csvfile.Close()

	return nil
}

func getAllStoresInSales(sales []sale) []string {
	stores := make([]string, 0)

	for _, s := range sales {
		storeExists := false
		for _, store := range stores {
			if store == s.storeName {
				storeExists = true
			}
		}
		if !storeExists {
			stores = append(stores, s.storeName)
		}
	}

	return stores
}

func getAllStoreSales(sales []sale, store string) []sale {
	storeSales := make([]sale, 0)

	for _, s := range sales {
		if s.storeName == store {
			storeSales = append(storeSales, s)
		}
	}
	return storeSales

}

func MakePurchasesCSVs(sales []sale) error {

	stores := getAllStoresInSales(sales)

	//now we're going to loop through the stores, and make a file for each one
	for _, store := range stores {
		//we need to make a file-name-safe value from the store name
		filename := strings.ToLower(store)

		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		filename = fmt.Sprintf("%s.csv", reg.ReplaceAllString(filename, ""))

		csvfile, err := os.Create(filename)

		if err != nil {
			return err
		}

		csvwriter := altcsv.NewWriter(csvfile)
		csvwriter.AllQuotes = true

		headerRow := make([]string, 0)
		headerRow = append(headerRow, "invoice date")
		headerRow = append(headerRow, "invoice number")
		headerRow = append(headerRow, "din")
		headerRow = append(headerRow, "qty")

		err = csvwriter.Write(headerRow)
		if err != nil {
			return err
		}

		storeSales := getAllStoreSales(sales, store)

		for _, s := range storeSales {
			dataRow := make([]string, 0)
			dataRow = append(dataRow, s.saleDate)
			dataRow = append(dataRow, s.rxNum)
			dataRow = append(dataRow, strconv.Itoa(s.din))
			dataRow = append(dataRow, strconv.FormatFloat(s.quantity, 'f', 1, 64))
			err = csvwriter.Write(dataRow)
			if err != nil {
				return err
			}
		}
		csvwriter.Flush()

		csvfile.Close()
		fmt.Println(filename, " written")
	}
	return nil
}
