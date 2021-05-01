package main

import (
	"citation-graph/backend/graph"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type S2Record struct {
	ID    string `json:"paper_id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	// Abstract string `json:"abstract"`
	InCitations  []string `json:"inbound_citations"`  // Papers that ite this paper
	OutCitations []string `json:"outbound_citations"` // Papers that ite this paper
}

func main() {
	ctx := context.Background()
	corpusRoot := "/home/jordan/s2orc-2020-07-05-v1/full/metadata/"
	badgerGraph, err := graph.NewBadgerCitationGraph("./badger")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	badgerGraph.GenerateGraphFromCorpus(ctx, corpusRoot)

	count, err := badgerGraph.CountNodes()
	fmt.Println(count)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
