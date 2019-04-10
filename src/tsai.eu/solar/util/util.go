package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"errors"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

//------------------------------------------------------------------------------

// ConvertYAMLToJSON converts YAML to JSON
func ConvertYAMLToJSON(input []byte) ([]byte, error) {
	var yamlData interface{}
	var jsonData interface{}
	var output   []byte
	
	err := yaml.Unmarshal(input, &yamlData)
	if err == nil {
		jsonData, err = convertYAMLToJSON(yamlData)
	}
	if err == nil {
		output, err = json.Marshal(jsonData)
	}
	return output, err
}

func convertYAMLToJSON(inputObj interface{}) (interface{}, error) {
	switch inputObj.(type) {
	case map[interface{}]interface{}:
		input := inputObj.(map[interface{}]interface{})
		output := make(map[string]interface{})
		for key, value := range input {
			var outputKey string
			switch key.(type) {
			case string:
				outputKey = key.(string)
			case int:
				outputKey = strconv.Itoa(key.(int))
			default:
				return nil, fmt.Errorf("Expected map key to be a string or int, but was %T", key)
			}
			outputValue, err := convertYAMLToJSON(value)
			if err != nil {
				return nil, err
			}
			output[outputKey] = outputValue
		}
		return output, nil
	case []interface{}:
		input := inputObj.([]interface{})
		output := make([]interface{}, len(input))
		for i, inputElem := range input {
			outputElem, err := convertYAMLToJSON(inputElem)
			if err != nil {
				return nil, err
			}
			output[i] = outputElem
		}
		return output, nil
	default:
		return inputObj, nil
	}
}

//------------------------------------------------------------------------------

// LoadFile reads string from a file
func LoadFile(filename string) (data string, err error) {
	bytes, err := ioutil.ReadFile(filename)

	data = (string)(bytes)

	return
}

//------------------------------------------------------------------------------

// SaveFile writes string to a file
func SaveFile(filename string, data string) (err error) {
	err = ioutil.WriteFile(filename, []byte(data), 0644)

	return
}

//------------------------------------------------------------------------------

// LoadYAML reads yaml from a file and transforms into the structure of the entity
func LoadYAML(filename string, entity interface{}) (err error) {
	// handle panic errors
	// defer func() {
	// 	if r := recover(); r!= nil {
	// 		fmt.Println(r)
	// 		err = errors.New("invalid structure")
	// 	}
	// }()

	// read file
	yamlbytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return errors.New("unable to load data: " + err.Error())
	}

	err = errors.New("invalid structure")

	// // unmarshal data
	err = yaml.Unmarshal(yamlbytes, entity)

	if err != nil {
		return errors.New("invalid structure: " + err.Error())
	}

	// success
	return err
}

//------------------------------------------------------------------------------

// SaveYAML writes the entity as yaml data to a file
func SaveYAML(filename string, entity interface{}) error {
	// marshal entity
	bytes, err := yaml.Marshal(entity)

	// write the entity
	err = ioutil.WriteFile(filename, bytes, 0644)

	if err != nil {
		return errors.New("unable to save entity: " + err.Error())
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// DumpYAML writes the entity as yaml data to the console
func DumpYAML(entity interface{}) {
	// marshal entity
	bytes, _ := yaml.Marshal(entity)

	fmt.Println(string(bytes))
}

//------------------------------------------------------------------------------

// ConvertFromYAML transforms yaml into the structure of the entity
func ConvertFromYAML(yaml string, entity interface{}) error {
	// convert to JSON
	jsonbytes, err := ConvertYAMLToJSON([]byte(yaml))

	// unmarshal data
	err = json.Unmarshal(jsonbytes, entity)

	if err != nil {
		return errors.New("invalid structure: " + err.Error())
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// ConvertToYAML show entity as yaml
func ConvertToYAML(entity interface{}) (string, error) {
	// marshal entity
	bytes, err := yaml.Marshal(entity)

	// success
	return string(bytes), err
}

//------------------------------------------------------------------------------

// UUID creates a universal unique id
func UUID() string {
	return uuid.New().String()
}

//------------------------------------------------------------------------------

// Print outputs a message
func Print(format string, args ...interface{} )  {
	fmt.Printf(format, args...)
}

//------------------------------------------------------------------------------
