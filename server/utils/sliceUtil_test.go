package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddIds(t *testing.T) {
	// Test adding a new key to an empty slice
	ids := []string{}
	key := "123"
	expected := []string{"123"}
	result := AppendSlice(ids, key)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test adding a key that already exists in the slice
	ids = []string{"123", "456"}
	key = "123"
	expected = []string{"123", "456"}
	result = AppendSlice(ids, key)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test adding a new key to a non-empty slice
	ids = []string{"123", "456"}
	key = "789"
	expected = []string{"123", "456", "789"}
	result = AppendSlice(ids, key)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestAddJsonIds(t *testing.T) {
	// Test adding a new key to an empty JSON array
	jsonArr := ""
	key := "123"
	expected := `["123"]`
	result := AppendJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test adding a key that already exists in the JSON array
	jsonArr = `["123","456"]`
	key = "123"
	expected = `["123","456"]`
	result = AppendJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test adding a new key to a non-empty JSON array
	jsonArr = `["123","456"]`
	key = "789"
	expected = `["123","456","789"]`
	result = AppendJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test adding an invalid JSON array
	jsonArr = "invalid"
	key = "123"
	expected = "invalid"
	result = AppendJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestIsExistSlice(t *testing.T) {
	// Test a key that exists in the slice
	ids := []string{"123", "456"}
	key := "123"
	expected := true
	result := IsExistSlice(ids, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test a key that does not exist in the slice
	ids = []string{"123", "456"}
	key = "789"
	expected = false
	result = IsExistSlice(ids, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestIsExistJsonSlice(t *testing.T) {
	// Test a key that exists in the JSON array
	jsonArr := `["123","456"]`
	key := "123"
	expected := true
	result := IsExistJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test a key that does not exist in the JSON array
	jsonArr = `["123","456"]`
	key = "789"
	expected = false
	result = IsExistJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test an invalid JSON array
	jsonArr = "invalid"
	key = "123"
	expected = false
	result = IsExistJsonSlice(jsonArr, key)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestPrintCharactersWithInterval(t *testing.T) {
	data := "hello, \n你好"

	fmt.Println(len(data))
	for _, c := range data {
		fmt.Print(string(c))
	}

	fmt.Println()

	for i := 0; i < len(data); i++ {
		fmt.Print(string(data[i]))
	}
}
