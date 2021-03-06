package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
)

type S2Record struct {
	ID          string   `json:"id"`
	InCitations []string `json:"inCitations"`
}

func main() {
	fmt.Println("hello!")
	f, err := os.Open("../s2-corpus-994.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	defer gr.Close()
	scanner := bufio.NewScanner(gr)
	for scanner.Scan() {
		record := &S2Record{}
		json.Unmarshal(scanner.Bytes(), record)
		fmt.Println(record)
	}
}
