package generators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/gengo/types"
)

// Validation Markers.
const (
	Maximum          = "Maximum"
	Minimum          = "Minimum"
	ExclusiveMaximum = "ExclusiveMaximum"
	ExclusiveMinimum = "ExclusiveMinimum"
	MaxLength        = "MaxLength"
	MinLength        = "MinLength"
	MaxProperties    = "MaxProperties"
	MinProperties    = "MinProperties"
	MaxItems         = "MaxItems"
	MinItems         = "MinItems"
	Pattern          = "Pattern"
	UniqueItems      = "UniqueItems"
)

type validationParser struct {
	generator  openAPITypeWriter
	member     *types.Member
	typeString string
}

func GenerateValidationProperty(vp validationParser) error {
	for _, val := range vp.member.CommentLines {
		ruleMap := make(map[string]string)
		if err := splitMarker(val, ruleMap); ruleMap != nil && err == nil {
			for ruleKey, ruleValue := range ruleMap {
				if err := applyRule(vp, ruleKey, ruleValue); err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func validationParseError(ruleKey, ruleValue string, m *types.Member) error {
	return fmt.Errorf("failed to parse validation property %v=%v for field %v", ruleKey, ruleValue, m)
}

func applyRule(vp validationParser, ruleKey, ruleValue string) error {
	switch ruleKey {

	// string markers
	case MaxLength, MinLength:
		if vp.typeString != "string" {
			return fmt.Errorf("must apply maxlength, minlength to a string, found %s", vp.typeString)
		}
		i, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		val := &i
		vp.generator.Do(fmt.Sprintf("%s: IntPtr($.$),\n", ruleKey), val)
	case Pattern:
		if vp.typeString != "string" {
			return fmt.Errorf("must apply pattern to a string, found %s", vp.typeString)
		}
		vp.generator.Do(fmt.Sprintf("%s: \"$.$\",\n", ruleKey), ruleValue)

	//numeric markers
	case Maximum, Minimum:
		if !isNumericType(vp.typeString) {
			return fmt.Errorf("must apply maximum, minimum to a numeric value, found %s", vp.typeString)
		}
		f, err := strconv.ParseFloat(ruleValue, 64)
		if err != nil {
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		val := &f
		vp.generator.Do(fmt.Sprintf("%s: FloatPtr($.$),\n", ruleKey), val)
	case ExclusiveMaximum, ExclusiveMinimum:
		if !isNumericType(vp.typeString) {
			return fmt.Errorf("must apply exclusivemaximum, exclusiveminimum to a numeric value, found %s", vp.typeString)
		}
		b, err := strconv.ParseBool(ruleValue)
		if err != nil {
			fmt.Println(validationParseError(ruleKey, ruleValue, vp.member))
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		vp.generator.Do(fmt.Sprintf("%s: $.$,\n", ruleKey), b)

	case MaxProperties, MinProperties:
		if !isNumericType(vp.typeString) {
			return fmt.Errorf("must apply maxproperties, minproperties to a numeric value, found %s", vp.typeString)
		}
		i, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		val := &i
		vp.generator.Do(fmt.Sprintf("%s: IntPtr($.$),\n", ruleKey), val)

	// slice markers
	case MaxItems, MinItems:
		if vp.typeString != "array" {
			return fmt.Errorf("must apply maxitem, minitem to an array, found %s", vp.typeString)
		}
		i, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		val := &i
		vp.generator.Do(fmt.Sprintf("%s: IntPtr($.$),\n", ruleKey), val)
	case UniqueItems:
		if vp.typeString != "array" {
			return fmt.Errorf("must apply uniqueitems to an array, found %s", vp.typeString)
		}
		b, err := strconv.ParseBool(ruleValue)
		if err != nil {
			return validationParseError(ruleKey, ruleValue, vp.member)
		}
		vp.generator.Do(fmt.Sprintf("%s: $.$,\n", ruleKey), b)

	default:
		return fmt.Errorf("unsupported validation rule <%v> for field <%v>", ruleKey, vp.member)
	}
	return nil
}

func splitMarker(raw string, ruleMap map[string]string) (err error) {
	// Ignore all the lines with no validation markers.
	if !strings.HasPrefix(raw, "nexus-validation") {
		return nil
	}

	re := regexp.MustCompile(`^(nexus-validation:)|\s`)
	rules := re.ReplaceAllString(raw, "")

	rulesParts := strings.Split(rules, ",")
	for _, rule := range rulesParts {
		parts := strings.Split(rule, "=")

		if len(parts) < 2 {
			return fmt.Errorf(`please check if the validation rule is added in the format ruleName:ruleValue` +
				`for example: nexus-validation: MaxLength=8`)
		}
		ruleMap[parts[0]] = parts[1]
	}
	return nil
}

func isNumericType(typeString string) bool {
	return typeString == "integer" || typeString == "number"
}
