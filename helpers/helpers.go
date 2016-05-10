package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
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

func JSONStringToStructure(jsonString string, structure interface{}) error {
	jsonBytes := []byte(jsonString)
	return json.Unmarshal(jsonBytes, structure)
}

func StructureToJSON(structure interface{}) (string, error) {
	bin, err := json.Marshal(structure)
	return string(bin), err
}

// read content from file
// path containing a dash will mean read from std in
func ReadFromFile(path string) (string, error) {
	// check for a dash
	if len(path) == 1 && path == "-" {
		// read from input
		var buffer bytes.Buffer
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return buffer.String(), nil
	} else if len(path) > 1 {
		// read file
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	}
	return "", nil
}
