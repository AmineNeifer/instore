package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/jwangsadinata/go-multimap"
	"github.com/jwangsadinata/go-multimap/setmultimap"
)

const (
	bCsv = "[csv] "
)
// fromMultiToSlice is a function that gets a multimap
// and converts it to a 2D slice of strings
func fromMultiToSlice(m multimap.MultiMap) [][]string {
	var records [][]string
	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			record := []string{k.(string), v.(string)}
			records = append(records, record)
		}
	}
	return records
}

// SaveCsv is a fuction that saves multimap into a csv file
func SaveCsv(m multimap.MultiMap, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf(bCsv+"Failed to create file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	records := fromMultiToSlice(m)

	w := csv.NewWriter(f)

	err = w.WriteAll(records) // calls Flush internally
	if err != nil {
		fmt.Printf(bCsv+"Failed to write to file: %v", err)
		os.Exit(1)
	}
}

// fromMultiToSlice is a function that gets a 2D slice of strings
// and converts it to a setmultimap
func fromSliceToMulti(records [][]string) multimap.MultiMap {
	m := setmultimap.New()

	for _, record := range records {
		m.Put(record[0], record[1])
	}
	return m
}

// LoadCsv is a function that loads csv file into a multimap
func LoadCsv(filename string) multimap.MultiMap {
	// if file is doesn't exist it returns an empty setmultimap
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return setmultimap.New()
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf(bCsv+"Failed to open file: %v", err)
		os.Exit(1)
	}

	defer f.Close()	

	r := csv.NewReader(f)

	// Get first line
	record, err := r.Read()

	// if file is empty it returns an empty setmultimap
	if err == io.EOF {
		return setmultimap.New()
	}

	records, err := r.ReadAll()

	// add first line (record) to records
	records = append(records, record)
	if err != nil {
		fmt.Printf(bCsv+"Failed to read file: %v", err)
		os.Exit(1)
	}

	return fromSliceToMulti(records)
}
