package util

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

//func ToPlural(singular string) string {
//	plural := namer.NewAllLowercasePluralNamer(nil)
//	t := &types.Type{
//		Name: types.Name{
//			Name: singular,
//		},
//	}
//	return plural.Name(t)
//}

func RemoveSpecialChars(value string) string {
	re, err := regexp.Compile(`[\_\.\/\-]`)
	if err != nil {
		log.Fatalf("failed to remove special chars from string %v: %v", value, err)
	}
	return re.ReplaceAllString(value, "")
}
