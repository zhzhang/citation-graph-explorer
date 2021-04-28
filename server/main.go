package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

type S2Record struct {
	ID          string   `json:"paper_id"` 
	Title    string `json:"title"`
	Year 	int `json:"year"`
	// Abstract string `json:"abstract"`
	InCitations []string `json:"inbound_citations"` // Papers that ite this paper
	OutCitations []string `json:"outbound_citations"` // Papers that ite this paper
}

func main() {
	// ctx := context.Background()
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	count := 0
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	  })
	fmt.Println(count)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}