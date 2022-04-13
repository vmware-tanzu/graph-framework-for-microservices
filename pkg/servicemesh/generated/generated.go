package generated

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type GeneratedObjectInterface interface {
	GetUrl(map[string]string) (string, error)
	GetBody(interface{}) ([]byte, error)
}

func GetGeneratedMapKey(keyList []string) (key string) {
	for _, k := range keyList {
		if len(key) == 0 {
			key = k
		} else {
			key = key + ":" + k
		}
	}
	return key
}

var GeneratedMap map[string]GeneratedObjectInterface

func ApplyDecode(inputMetadata yaml.MapSlice, spec interface{}) (url string, body []byte, err error) {

	var objKey []string
	metadata := make(map[string]string)
	for _, value := range inputMetadata {
		strKey := fmt.Sprintf("%v", value.Key)
		strValue := fmt.Sprintf("%v", value.Value)
		metadata[strKey] = strValue
		objKey = append(objKey, strKey)

	}

	fmt.Printf("objKey: %v\n", objKey)
	fmt.Printf("metadata: %v\n", metadata)
	fmt.Printf("GetGeneratedMapKey: %v\n", GetGeneratedMapKey(objKey))
	fmt.Printf("GeneratedMap: %v\n", GeneratedMap)

	if obj, ok := GeneratedMap[GetGeneratedMapKey(objKey)]; ok {
		url, err := obj.GetUrl(metadata)
		fmt.Printf("Generated key: %v, %v\n", url, err)

		body, err := obj.GetBody(spec)
		fmt.Printf("Returned JSON body: %+v", body)
		return url, body, nil
	}

	return "", nil, errors.New(fmt.Sprintf("unable to decode apply for input"))
}

func DeleteDecode(inputMetadata yaml.MapSlice) (url string, err error) {

	var objKey []string
	metadata := make(map[string]string)
	for _, value := range inputMetadata {
		strKey := fmt.Sprintf("%v", value.Key)
		strValue := fmt.Sprintf("%v", value.Value)
		metadata[strKey] = strValue
		objKey = append(objKey, strKey)

	}

	fmt.Printf("objKey: %v\n", objKey)
	fmt.Printf("metadata: %v\n", metadata)
	fmt.Printf("GetGeneratedMapKey: %v\n", GetGeneratedMapKey(objKey))
	fmt.Printf("GeneratedMap: %v\n", GeneratedMap)

	if obj, ok := GeneratedMap[GetGeneratedMapKey(objKey)]; ok {
		url, err := obj.GetUrl(metadata)
		fmt.Printf("Generated key: %v, %v\n", url, err)

		return url, nil
	}

	return "", errors.New(fmt.Sprintf("unable to decode delete for input"))
}

func init() {
	GeneratedMap = make(map[string]GeneratedObjectInterface)
}
