package synmedreader

import (
	"errors"
	"path/filepath"
)

func ProcessFile(fName string) ([]sale, error) {
	extension := filepath.Ext(fName)
	if extension == ".csv" {
		sales, err := processCSV(fName)
		if err != nil {
			return make([]sale, 0), err
		}
		return sales, nil
	} else if extension == ".xls" {
		sales, err := processXLS(fName)
		if err != nil {
			return make([]sale, 0), err
		}
		return sales, nil
	}
	return make([]sale, 0), errors.New("Unknown error occurred during file processing")
}
