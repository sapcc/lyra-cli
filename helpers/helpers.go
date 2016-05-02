package helpers

import (
	"regexp"
	"strings"
)

func StringTokeyValueMap(data string) map[string]string {
	result_map := make(map[string]string)
	if len(data) > 0 {
		keyValues := strings.Split(data, ",")
		for _, kv := range keyValues {
			reg := regexp.MustCompile(`\:|\=`)
			tags_array := reg.Split(kv, -1)
			if len(tags_array) == 2 {
				result_map[tags_array[0]] = tags_array[1]
			}
		}
	} else {
		return map[string]string{}
	}

	return result_map
}

func StringToArray(data string) []string {
	result_array := []string{}
	if len(data) > 0 {
		keyValues := strings.Split(data, ",")
		for _, value := range keyValues {
			result_array = append(result_array, value)
		}
	} else {
		return []string{}
	}
	return result_array
}
