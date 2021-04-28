package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	badger "github.com/dgraph-io/badger/v3"
	"golang.org/x/sync/errgroup"
)

type S2Record struct {
	ID          string   `json:"paper_id"` 
	Title    string `json:"title"`
	Year 	int `json:"year"`
	// Abstract string `json:"abstract"`
	InCitations []string `json:"inbound_citations"` // Papers that cite this paper.
	OutCitations []string `json:"outbound_citations"` // Papers that this paper cites.
	DDI float64 // Discounted downstream impact.
}

func main() {
	ctx := context.Background()
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loadGraph(ctx, db)
}

func loadGraph(ctx context.Context, db *badger.DB) {
	corpusRoot := "/home/jordan/s2orc-2020-07-05-v1/full/metadata/"
	files, err := ioutil.ReadDir(corpusRoot)
	if err != nil {
		log.Fatal(err)
		return
	}
	countChannel := make(chan int)
	eg, _ := errgroup.WithContext(ctx)
	for _, f := range files {
		go func(filename string) {
			eg.Go(func() error {
				return processFile(filename, db, countChannel)
			})
		}(path.Join(corpusRoot, f.Name()))
	}
	total := 0
	go func() {
		for v := range countChannel {
			total += v
		}
		fmt.Println("done")
	}()
	eg.Wait()
	close(countChannel)
	fmt.Println(total)
}


func processFile(file string, db *badger.DB, c chan int) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	scanner := bufio.NewScanner(gr)
	count := 0
	for scanner.Scan() {
		count++
		record := S2Record{}
		json.Unmarshal(scanner.Bytes(), &record)
		db.Update(func(txn *badger.Txn) error {
			var b bytes.Buffer
			e := gob.NewEncoder(&b)
			if err := e.Encode(record); err != nil {
				panic(err)
			}
			err := txn.Set([]byte(record.ID), b.Bytes())
			return err
		})
	}
	c <- count
	log.Printf("finished processing %s", file)
	return nil
}