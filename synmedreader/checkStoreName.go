package synmedreader

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

func checkStoreName(storeName string) (string, error) {
	options, err := getNameOptions(storeName)
	if err != nil {
		return storeName, err
	}

	selection, err := selectName(options, storeName)
	if err != nil {
		return storeName, err
	}

	//We start prompt new name with default value of current storeName
	if selection == "Custom" {
		selection, err = createName(storeName)
		if err != nil {
			return storeName, err
		}
	}

	return selection, nil
}

func getNameOptions(storeName string) ([]string, error) {
	options := make([]string, 0)
	if strings.Contains(storeName, "-") {
		//This means that there is a dash in the store Name
		// Perform the split, clean the strings, and add them up
		storeNameParts := strings.Split(storeName, "-")
		for i := 0; i < len(storeNameParts); i++ {
			storeNamePart, err := cleanString(storeNameParts[i])
			if err == nil && len(storeNamePart) > 0 {
				if strings.Contains(storeNamePart, " ") {
					options = append(options, storeNamePart)
				} else {
					options = append(options, "Store "+storeNamePart)
				}
			}
		}
	}

	fullStoreName, err := cleanString(storeName)
	if err == nil && len(fullStoreName) > 0 {
		if strings.Contains(fullStoreName, " ") {
			options = append(options, fullStoreName)
		} else {
			options = append(options, "Store "+fullStoreName)
		}
	}

	return append(options, "Custom"), nil
}

func selectName(nameOptions []string, storeName string) (string, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Store Found: '%s', How should this appear in the exported file?", storeName),
		Items: nameOptions,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return storeName, err
	}

	return result, nil
}

func createName(storeName string) (string, error) {
	validate := func(input string) error {
		input, cleanErr := cleanString(input)
		if cleanErr != nil {
			return cleanErr
		}
		if len(input) < 2 {
			return errors.New("Please enter a valid store name")
		}
		if !strings.Contains(input, " ") {
			return errors.New("Please enter a name that has 2 words in it (to represent first/last patient name)")
		}
		return nil
	}

	defaultStoreName, err := cleanString(storeName)
	if err != nil {
		defaultStoreName = ""
	}

	prompt := promptui.Prompt{
		Label:    "Please enter a store name",
		Validate: validate,
		Default:  defaultStoreName,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func cleanString(str string) (string, error) {
	replacedPoundWithSpace := strings.Replace(str, "#", " ", -1)

	//Remove all values that are not (alpha/numeric/spaces)
	regSpecialChars, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		return "", err
	}
	removedSpecialCharacters := regSpecialChars.ReplaceAllString(replacedPoundWithSpace, "")

	//Remove all excess whitespace
	regExtraWhiteSpace, err := regexp.Compile(`\s+`)
	if err != nil {
		return "", err
	}
	removedDuplicateWhiteSpace := regExtraWhiteSpace.ReplaceAllString(removedSpecialCharacters, " ")

	return strings.TrimSpace(removedDuplicateWhiteSpace), nil

}
