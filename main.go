package main

import (
	"fmt"
	"log"
	"synmed-reader/synmedreader"
)

func main() {
	//Lets start by getting the files in the current directory that could be the file
	files, err := synmedreader.RetrieveFiles()
	if err != nil {
		log.Fatal(err)
		return
	}

	//Prompt or confirm with the user the file they want to process
	selectedFile, err := synmedreader.SelectFile(files)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Now we need to process the file - we call the Process File function, which will automatically select the best processing depending on the file extension
	sales, err := synmedreader.ProcessFile(selectedFile)

	//Now that we have the file, let's output it to CSV
	if err = synmedreader.MakeCSV(sales); err != nil {
		log.Fatal(err)
		return
	}

	//This is where we loop through and create the purchase records for the satellite stores
	if err = synmedreader.MakePurchasesCSVs(sales); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("CentralSales.csv has been written")
}
