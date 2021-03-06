package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type S2Record struct {
	ID          string   `json:"id"`
	InCitations []string `json:"inCitations"`
}

func main() {
	f, err := os.Open("../s2-corpus-1994.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	defer gr.Close()
	scanner := bufio.NewScanner(gr)
	t := time.Now()
	count := 0
	for scanner.Scan() {
		count++
		record := &S2Record{}
		json.Unmarshal(scanner.Bytes(), record)
		fmt.Println(record)
	}
	fmt.Println(time.Since(t))
	fmt.Println(count)
}
