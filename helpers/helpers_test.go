package helpers

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStringTokeyValueMap(t *testing.T) {
	result := StringTokeyValueMap("")
	eq := reflect.DeepEqual(result, map[string]string{})
	if !eq {
		t.Error("Expected StringTokeyValueMap to be return empty map")
	}
	result = StringTokeyValueMap("name:test")
	eq = reflect.DeepEqual(result, map[string]string{"name": "test"})
	if !eq {
		t.Error(`Expected StringTokeyValueMap to be map[string]string{"name": "test"}`)
	}
	result = StringTokeyValueMap("name:test,test1=test1")
	eq = reflect.DeepEqual(result, map[string]string{"name": "test", "test1": "test1"})
	if !eq {
		t.Error(`Expected StringTokeyValueMap to be map[string]string{"name": "test", "test1": "test1"}`)
	}
	result = StringTokeyValueMap("name:test,test1=test1,test1:miau")
	eq = reflect.DeepEqual(result, map[string]string{"name": "test", "test1": "miau"})
	if !eq {
		t.Error(`Expected StringTokeyValueMap to be map[string]string{"name": "test", "test1": "miau"}`)
	}
	result = StringTokeyValueMap("name:test,test1,")
	eq = reflect.DeepEqual(result, map[string]string{"name": "test"})
	if !eq {
		t.Error(`Expected ignore broken key pairs`)
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
