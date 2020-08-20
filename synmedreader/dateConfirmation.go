package synmedreader

import (
	"errors"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

func confirmDate(dts string) (string, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Discovered report end date of %s, would you like to use this as the sale date?", dts),
		IsConfirm: true,
	}

	_, err := prompt.Run()

	//I think this means they say no
	if err == nil {
		return dts, nil

	}

	newDate, err := getNewDate()
	if err == nil {
		return newDate, nil
	}

	return "", err
}

func getNewDate() (string, error) {

	validate := func(input string) error {
		_, err := time.Parse("2006-01-02", input)
		if err != nil {
			return errors.New("Make sure the date is in the format of YYYY-MM-DD")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Ok - great, what date should the report have? (YYYY-MM-DD)",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}
