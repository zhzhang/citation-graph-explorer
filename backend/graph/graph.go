package graph

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

type Record struct {
	ID           string   `json:"paper_id"`
	Title        string   `json:"title"`
	Year         int      `json:"year"`
	InCitations  []string `json:"inbound_citations"`  // Papers that cite this paper.
	OutCitations []string `json:"outbound_citations"` // Papers that this paper cites.
	DDI          float64  // Discounted downstream impact.
	// Abstract string `json:"abstract"`
}

type CitationGraph interface {
	GetNode(id string) Record
	GenerateGraphFromCorpus(corpusPath string)
	CountNodes() (int, error)
}

type BadgerCitationGraph struct {
	db *badger.DB
}

func NewBadgerCitationGraph(dbPath string) (*BadgerCitationGraph, error) {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}
	return &BadgerCitationGraph{
		db: db,
	}, nil
}

func (g *BadgerCitationGraph) GetNode(id string) *Record {
	return &Record{}
}

func (g *BadgerCitationGraph) CountNodes() (int, error) {
	count := 0
	err := g.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (g *BadgerCitationGraph) GenerateGraphFromCorpus(ctx context.Context, corpusPath string) {
	files, err := ioutil.ReadDir(corpusPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	countChannel := make(chan int)
	eg, _ := errgroup.WithContext(ctx)
	for _, f := range files {
		go func(filename string) {
			eg.Go(func() error {
				return g.processFile(filename, countChannel)
			})
		}(path.Join(corpusPath, f.Name()))
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

func (g *BadgerCitationGraph) processFile(file string, c chan int) error {
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
		record := Record{}
		json.Unmarshal(scanner.Bytes(), &record)
		g.db.Update(func(txn *badger.Txn) error {
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
