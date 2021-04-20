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
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type S2Record struct {
	ID          string   `json:"id"` 
	InCitations []string `json:"inCitations"` // Papers that ite this paper
	OutCitations []string `json:"outCitations"` // Papers that ite this paper
}

func main() {
	// ctx := context.Background()
	// db, err := badger.Open(badger.DefaultOptions("./badger"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// loadGraph(ctx, db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func loadGraph(ctx context.Context, db *badger.DB) {
	corpusRoot := "/home/jordan/s2corpus"
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
	}()
	eg.Wait()
	close(countChannel)
	fmt.Println(total)
}


func processFile(file string, db *badger.DB, c chan int) error {
	fmt.Println(file)
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
	t := time.Now()
	count := 0
	for scanner.Scan() {
		count++
		record := S2Record{}
		json.Unmarshal(scanner.Bytes(), &record)
		go func(record S2Record) {
			db.Update(func(txn *badger.Txn) error {
				var b bytes.Buffer
				e := gob.NewEncoder(&b)
				if err := e.Encode(record); err != nil {
					panic(err)
				}
				err := txn.Set([]byte(record.ID), b.Bytes())
				return err
			})
		}(record)
	}
	fmt.Println(time.Since(t))
	fmt.Println(count)
	c <- count
	return nil
}