package util

import (
	"math/rand"
	"strconv"
	"time"
)

// IntInStringSlice Judge if a string needle is in the integer candidates
func IntInStringSlice(needle string, candidates []int) bool {
	for _, candidate := range candidates {
		if needle == strconv.Itoa(candidate) {
			return true
		}
	}

	return false
}

// SliceUniqueStrings ...
func SliceUniqueStrings(slice []string) []string {
	hash := make(map[string]interface{})
	for _, value := range slice {
		hash[value] = nil
	}

	result := make([]string, len(hash))

	index := 0
	for value := range hash {
		result[index] = value
		index++
	}

	return result
}

// RandomSlice
//  @param list
//  @return []map
func RandomSlice(list []map[string]interface{}) []map[string]interface{} {

	if len(list) <= 1 {
		return list
	}

	//乱序
	rand.Seed(time.Now().UnixNano())
	for i := len(list) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		list[i], list[num] = list[num], list[i]
	}
	return list
}
