package synmedreader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// RetrieveFiles will return a list of all the files matching a CSV or XLS in the current working directory and returns a slice of strings with their names
func RetrieveFiles() ([]string, error) {
	dirname := "."
	f, err := os.Open(dirname)
	if err != nil {
		return make([]string, 0), err
	}
	files, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return make([]string, 0), err
	}

	items := make([]string, 0)

	for i := 0; i < len(files); i++ {
		extension := filepath.Ext(files[i])
		if files[i][0:1] != "." {
			if extension == ".csv" || extension == ".xls" {
				items = append(items, files[i])
			}
		}
	}

	if len(items) < 1 {
		return items, errors.New("No XLS or CSV files found")
	}

	return items, nil
}

//SelectFile will take a slice of strings for the file options and returns the selected file, and an error
func SelectFile(files []string) (string, error) {

	if len(files) == 1 {
		useFile, err := confirmFile(files[0])
		if err != nil {
			return "", err
		}

		if useFile {
			return files[0], nil
		}
	} else if len(files) > 1 {
		selectedItem, err := selectFileFromList(files)
		if err != nil {
			return "", err
		}
		return selectedItem, nil
	}
	return "", errors.New("No file selected")
}

func confirmFile(file string) (bool, error) {

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("One file found, Would you like to process this file?  %s", file),
		IsConfirm: true,
	}

	_, err := prompt.Run()

	if err != nil {
		return false, err
	}

	return true, nil
}

func selectFileFromList(files []string) (string, error) {

	prompt := promptui.Select{
		Label: "Please choose your file",
		Items: files,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}
