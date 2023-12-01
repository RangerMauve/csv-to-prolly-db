package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	indexer "github.com/RangerMauve/ipld-prolly-indexer/indexer"
	car "github.com/ipld/go-car/v2"
	carBlockstore "github.com/ipld/go-car/v2/blockstore"

	cid "github.com/ipfs/go-cid"
	dagjson "github.com/ipld/go-ipld-prime/codec/dagjson"
	datamodel "github.com/ipld/go-ipld-prime/datamodel"
	qp "github.com/ipld/go-ipld-prime/fluent/qp"
	basicnode "github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/urfave/cli/v2"
)

// This is the default CID that gets
const EMPTY_DB_ROOT = "bafyrefczokuljxpuzx3ivzun5p5jdfnfdj3qzqq"

func main() {
	var output string
	var input string
	var collection string

	app := &cli.App{
		Name:  "csv-to-prolly-db",
		Usage: "Ingest a CSV file into a prolly tree",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "./db.car",
				Usage:       "Path to save database CAR file to.",
				Destination: &output,
			},
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Value:       "",
				Usage:       "Path to CSV file to read. Omit to load from STDIN",
				Destination: &input,
			},
			&cli.StringFlag{
				Name:        "collection",
				Aliases:     []string{"c"},
				Value:       "default",
				Usage:       "Name of database collection to save into",
				Destination: &collection,
			},
		},
		Action: func(*cli.Context) error {
			var reader io.Reader

			if input == "" {
				reader = os.Stdin
			} else {
				file, err := os.Open(input)
				if err != nil {
					return err
				}

				defer file.Close()

				reader = file
			}

			return Run(output, reader, collection)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(
	output string,
	reader io.Reader,
	collectionName string,
) error {

	ctx := context.Background()

	defaultCid, err := cid.Decode(EMPTY_DB_ROOT)

	if err != nil {
		return err
	}

	// Set a default CID so the CAR header has space for final CID
	blockstore, err := carBlockstore.OpenReadWrite(output, []cid.Cid{defaultCid})

	if err != nil {
		return err
	}

	db, err := indexer.NewDatabaseFromBlockStore(ctx, blockstore)

	if err != nil {
		return err
	}

	collection, err := db.Collection(collectionName)

	if err != nil {
		return err
	}

	_, err = collection.CreateIndex(ctx, "index")

	err = ingestCSV(ctx, reader, collection)

	if err != nil {
		return err
	}

	err = db.ApplyChanges(ctx)

	if err != nil {
		return err
	}
	finalCid := db.RootCid()

	err = blockstore.Finalize()

	if err != nil {
		return err
	}

	err = car.ReplaceRootsInFile(output, []cid.Cid{finalCid})

	if err != nil {
		return err
	}

	fmt.Println(finalCid)

	return nil
}

func ingestCSV(ctx context.Context, source io.Reader, collection *indexer.Collection) error {
	reader := csv.NewReader(source)

	headers, err := reader.Read()

	numFields := int64(len(headers) + 1)

	if err != nil {
		return err
	}

	index := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		node, err := qp.BuildMap(basicnode.Prototype.Any, numFields, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, "index", qp.Int(int64(index)))
			for fieldIndex, fieldValue := range record {
				nb := basicnode.Prototype__Any{}.NewBuilder()
				err := dagjson.Decode(nb, strings.NewReader(fieldValue))

				// If it wasn't json, it's just a string
				if err != nil {
					qp.MapEntry(ma, headers[fieldIndex], qp.String(fieldValue))
				} else {
					value := nb.Build()
					qp.MapEntry(ma, headers[fieldIndex], qp.Node(value))
				}
			}
		})

		if err != nil {
			return err
		}

		err = collection.Insert(ctx, node)

		if err != nil {
			return err
		}

		index++
	}

	return nil
}
