package utils

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

type StringArray []string

// Contains returns true if the string array contains the element
func (s StringArray) Contains(e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return -1, err
	}

	return num, nil
}

func StringToUInt(str string) (uint, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return 0, err
	}

	return uint(num), nil
}
