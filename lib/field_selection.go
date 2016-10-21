package lib

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structs"

	"errors"
	"fmt"
	"reflect"
	"strings"
)

var PR_DUMP = spew.Dump
var PR_INFO = fmt.Println

func FieldSelection(v interface{}, field string) (map[string]interface{}, error) {
	var fields []string
	var err error
	map_out := make(map[string]interface{})

	// If fields is empty, just return the whole struct
	if field == "" {
		s := structs.New(v)
		s.TagName = "json"
		return s.Map(), nil
	}

	// Check the input field, is it in the right format
	if fields, err = ValidateStringField(field); err != nil {
		return nil, err
	}

	// Loop through the field
	for _, field := range fields {
		field_name, _ := GetFieldNameFromJson(v, field)

		// Check if the field in the struct
		field_value, err := GetField(v, field_name)
		if err != nil {
			return nil, err
		}

		// if json tag specify the field's name, use it
		field_name, err = GetFieldJsonName(v, field_name, field_value)
		if err != nil {
			return nil, err
		}

		// add it to the map output
		map_out[field_name] = field_value
	}

	return map_out, nil
}

func ValidateStringField(field string) ([]string, error) {
	// Remove space if any
	field = strings.Replace(field, " ", "", -1)

	// Split the field by comma
	fields := strings.Split(field, ",")

	return fields, nil
}

func GetFieldJsonName(v interface{}, field_name string, field_value interface{}) (string, error) {
	s := structs.New(v)
	f := s.Field(field_name)

	// Get the value of field's json tag value
	json_tag := f.Tag("json")

	// Split the tag value by comma
	json_fields := strings.Split(json_tag, ",")

	// If there is no value in json tag -> len return is 0, use FieldName instead
	if len(json_fields) == 0 {
		return field_name, nil
	}

	// If user wants to select unexported json field, returns error
	if json_fields[0] == "-" {
		return "", errors.New(field_name + " is not exported to struct " + s.Name() + "'s json output")
	}

	// If the tags is ",omitempty"
	// -> the slice return is : len = 2, elem_1 = "", elem_2 = "omitempty"
	// If the tag has value ex, "access_token,omitempty"
	// -> elem_1 = "access_token"
	if json_fields[0] != "" {
		return json_fields[0], nil
	}

	// The json tag hasn't specified the field's return name, use the field name
	return field_name, nil
}

func GetFieldNameFromJson(v interface{}, json_name string) (string, error) {
	s := structs.New(v)
	fs := s.Fields()
	for _, f := range fs {
		json_tag := f.Tag("json")
		json_fields := strings.Split(json_tag, ",")
		if len(json_fields) == 0 {
			continue
		}
		if json_fields[0] == json_name {
			return f.Name(), nil
		}
	}
	return json_name, nil
}

func GetField(v interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil, errors.New(field + " is not part of " + r.Type().String())
	}
	return f.Interface(), nil
}
