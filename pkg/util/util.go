package util

import (
	"fmt"
	"regexp"
)

func ToPlural(singular string) string {
	var plural string

	if singular[len(singular)-1:] == "s" {
		plural = fmt.Sprintf("%ses", singular)
	} else if singular[len(singular)-1:] == "y" {
		plural = fmt.Sprintf("%sies", singular[:len(singular)-1])
	} else {
		plural = fmt.Sprintf("%ss", singular)
	}

	return plural
}

func RemoveSpecialChars(value string) string {
	re, err := regexp.Compile(`[\_\.]`)
	if err != nil {
		fmt.Println(err)
	}
	return re.ReplaceAllString(value, "")
}
