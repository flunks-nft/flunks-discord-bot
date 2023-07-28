package utils

import "github.com/joho/godotenv"

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
