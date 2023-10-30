package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	indexer "github.com/RangerMauve/ipld-prolly-indexer/indexer"

	ipld "github.com/ipld/go-ipld-prime"
	datamodel "github.com/ipld/go-ipld-prime/datamodel"
	basicnode "github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/zeebo/assert"
	"testing"
)

const EXAMPLE_DB_FILE = "db.car"

func TestIngestAndSearch(t *testing.T) {
	ctx := context.Background()

	collectionName := "test"

	output := EXAMPLE_DB_FILE

	reader := strings.NewReader(`name, message
Alice,Hello World
Bob,Goodbye World
`)

	err := Run(
		output,
		reader,
		collectionName,
	)

	assert.NoError(t, err)

	// Set a default CID so the CAR header has space for final CID
	db, err := indexer.ImportFromFile(output)

	assert.NoError(t, err)

	collection, err := db.Collection(collectionName)

	assert.NoError(t, err)

	node, err := getByIndex(ctx, 1, collection)

	assert.NoError(t, err)

	name, err := node.LookupByString("name")

	assert.NoError(t, err)

	assert.True(t, datamodel.DeepEqual(name, basicnode.NewString("Bob")))

}

func getByIndex(ctx context.Context, index int, collection *indexer.Collection) (ipld.Node, error) {
	value := strconv.Itoa(index)
	query := indexer.Query{
		Equal: map[string]ipld.Node{
			"index": basicnode.NewString(value),
		},
	}

	results, err := collection.Search(ctx, query)

	if err != nil {
		return nil, err
	}

	record, ok := <-results

	if !ok {
		return nil, fmt.Errorf("Unable to find record at index " + value)
	}

	return record.Data, nil
}
