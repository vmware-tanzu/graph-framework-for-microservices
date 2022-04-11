package util

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
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
	re, err := regexp.Compile(`[\_\.\/]`)
	if err != nil {
		log.Fatalf("failed to remove special chars from string %v: %v", value, err)
	}
	return re.ReplaceAllString(value, "")
}
