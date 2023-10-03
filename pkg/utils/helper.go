package utils

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("ENV") == "production" {
		// skip loading .env file if in production
		return
	}
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
	if str == "" {
		return -1, nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return -1, err
	}

	return num, nil
}

func StringToUInt(str string) (uint, error) {
	if str == "" {
		return 0, nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return 0, err
	}

	return uint(num), nil
}

func RandomItem(slice interface{}) interface{} {
	// Convert slice to reflect.Value
	s := reflect.ValueOf(slice)

	// Check if the provided interface is indeed a slice
	if s.Kind() != reflect.Slice {
		panic("Expected a slice!")
	}

	length := s.Len()
	if length == 0 {
		return nil // or panic with "empty slice!"
	}

	// Retrieve a random item from the slice
	return s.Index(rand.Intn(length)).Interface()
}

func GetRandomKeyFromMap(m map[string]string) (key string) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys[rand.Intn(len(keys))]
}
