package main

import (
	"os"
	"strconv"
)

func getVariable(key string, defaultValue string) (result string) {
	result = os.Getenv(key)
	if len(result) > 0 {
		return result
	}
	return defaultValue
}

func getBooleanVariable(key string, defaultValue bool) (result bool) {
	value := getVariable(key, strconv.FormatBool(defaultValue))
	result, _ = strconv.ParseBool(value)
	return
}
