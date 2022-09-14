package nexus

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/elliotchance/orderedmap"
)

const DEFAULT_KEY = "default"

func ParseCRDLabels(parents []string, labels map[string]string) *orderedmap.OrderedMap {
	m := orderedmap.NewOrderedMap()
	for _, parent := range parents {
		if label, ok := labels[parent]; ok {
			m.Set(parent, label)
		} else {
			m.Set(parent, DEFAULT_KEY)
		}
	}

	return m
}

func GetHashedName(crdName string, parents []string, labels map[string]string, name string) string {
	orderedLabels := ParseCRDLabels(parents, labels)

	var output string
	for i, key := range orderedLabels.Keys() {
		value, _ := orderedLabels.Get(key)

		output += fmt.Sprintf("%s:%s", key, value)
		if i < orderedLabels.Len()-1 {
			output += "/"
		}
	}

	output += fmt.Sprintf("%s:%s", crdName, name)

	h := sha1.New()
	h.Write([]byte(output))
	return hex.EncodeToString(h.Sum(nil))
}
