package helpers

import (
	"reflect"
	"strings"
	"testing"
)

func TestJSONStringToStructure(t *testing.T) {
	var testStructure interface{}
	testString := `{"test":"test"}`
	err := JSONStringToStructure(testString, &testStructure)
	if err != nil {
		t.Error(`JSONStringToStructure expected to not get an error`)
	}
	mapStructure := testStructure.(map[string]interface{})
	if mapStructure["test"] != "test" {
		t.Error(`JSONStringToStructure expected right convertion to structure`)
	}
}

func TestStructureToJSON(t *testing.T) {
	testStructure := map[string]interface{}{"miau": "bup"}
	testString, err := StructureToJSON(testStructure)
	if err != nil {
		t.Error(`StructureToJSON expected to not get an error`)
	}
	if !strings.Contains(testString, `{"miau":"bup"}`) {
		t.Error(`Command create chef expected to have same repository revision'`)
	}
}

func TestStringSliceValueMap(t *testing.T) {

	table := []struct {
		input  []string
		output map[string]string
	}{
		{[]string{}, map[string]string{}},
		{[]string{"name:test"}, map[string]string{"name": "test"}},
		{[]string{"name=test"}, map[string]string{"name": "test"}},
		{[]string{"name=test:bla"}, map[string]string{"name": "test:bla"}},
		{[]string{"name:test=bla"}, map[string]string{"name": "test=bla"}},
		{[]string{"name=test=bla"}, map[string]string{"name": "test=bla"}},
		{[]string{"name=test", "name2:test2,bla"}, map[string]string{"name": "test", "name2": "test2,bla"}},
	}

	for i, testCase := range table {
		result, err := StringSliceKeyValueMap(testCase.input)
		if err != nil {
			t.Errorf("Case %d failed. No error expected. err=%s", i, err)
		}
		if !reflect.DeepEqual(result, testCase.output) {
			t.Errorf("Case %d, Expected %#v, Got %#v", i, testCase.output, result)
		}
	}

	_, err := StringSliceKeyValueMap([]string{"bla"})
	if err == nil {
		t.Error("Expected StringTokenValueMap to return an error")
	}

}
func TestStringToArray(t *testing.T) {
	result := StringToArray("")
	eq := reflect.DeepEqual(result, []string{})
	if !eq {
		t.Error("Expected TestStringToArray to return empty array")
	}
	result = StringToArray("test,test2")
	eq = reflect.DeepEqual(result, []string{"test", "test2"})
	if !eq {
		t.Error(`Expected TestStringToArray to return []string{"test", "test2"`)
	}
}
