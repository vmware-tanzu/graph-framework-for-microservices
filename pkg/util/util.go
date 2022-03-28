package util

import "fmt"

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
