package util

import (
	"encoding/csv"
	"io"
	"os"
)

// ReadCsv
//  @param fileName
//  @return []
func ReadCsv(fileName string) [][]string {
	list := [][]string{}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		list = append(list, record)
	}
	return list
}
