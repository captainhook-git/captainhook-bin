package configuration

import (
	"fmt"
	"github.com/captainhook-go/captainhook/io"
	"strconv"
)

type Options struct {
	options map[string]interface{}
}

// AsBool tries to convert an option value to a bool and return it
func (o Options) AsBool(option string, defaultValue bool) bool {
	var value, ok = o.options[option]
	if ok {
		switch value := value.(type) {
		case bool:
			return value
		case int:
		case int64:
		case float64:
			return value > 0
		case string:
			return io.AnswerToBool(value)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

// AsInt tries to convert an option value to an int and return it
func (o Options) AsInt(option string, defaultValue int) int {
	var value, ok = o.options[option]
	if ok {
		switch value := value.(type) {
		case int:
			return value
		case int64:
		case float64:
			return int(value)
		case string:
			i, err := strconv.Atoi(value)
			if err != nil {
				return defaultValue
			}
			return i
		default:
			return defaultValue
		}
	}
	return defaultValue
}

// AsString tries to convert an option value to a string and return it
func (o Options) AsString(option string, defaultValue string) string {
	var value, ok = o.options[option]
	if ok {
		switch value := value.(type) {
		case int:
		case int64:
			return fmt.Sprintf("%d", value)
		case float64:
			return fmt.Sprintf("%f", value)
		case string:
			return value
		default:
			return defaultValue
		}
	}
	return defaultValue
}

// AsSliceOfStrings tries to convert an option value to a slice of strings and return it
func (o Options) AsSliceOfStrings(option string) []string {
	var data, ok = o.options[option]
	var strings []string
	if ok {
		switch v := data.(type) {
		case []interface{}:
			for _, item := range v {
				// Use type assertion to convert each element to a string
				if str, ok := item.(string); ok {
					strings = append(strings, str)
				}
			}
		}
	}
	return strings
}

func (o Options) All() map[string]interface{} {
	return o.options
}

func NewOptions(someMap map[string]interface{}) *Options {
	return &Options{options: someMap}
}
