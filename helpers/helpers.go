package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func KeyValueSplit(data string) (string, string, error) {
	parts := regexp.MustCompile(`:|=`).Split(data, 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("%v is not a valid key value pair. split result: %#v", data, parts)
}

func StringSliceKeyValueMap(data []string) (map[string]string, error) {
	result := map[string]string{}
	for _, elem := range data {
		key, value, err := KeyValueSplit(elem)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}

	return result, nil
}

func StringToArray(data string) []string {
	result_array := []string{}
	if len(data) > 0 {
		keyValues := strings.Split(data, ",")
		result_array = append(result_array, keyValues...)
	} else {
		return []string{}
	}
	return result_array
}

func JSONStringToStructure(jsonString string, structure interface{}) error {
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, structure)
	if err != nil {
		return errors.New(fmt.Sprint("Invalid JSON:: got: ", jsonString, ". ", err))
	}
	return nil
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
			if _, err := buffer.WriteString(scanner.Text()); err != nil {
				return "", err
			}
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return buffer.String(), nil
	} else if len(path) > 1 {
		// read file
		dat, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return "", err
		}
		return string(dat), nil
	}
	return "", nil
}

func MapMerge(dst, src interface{}) {
	dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)

	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}

func CheckErrAndPrintToStdErr(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s", msg, err.Error())
	}
}
