# csv-to-prolly-db

Take a csv file and turn it into a prolly tree based database indexed by row number.

Based on the [ipld-prolly-indexer](https://github.com/RangerMauve/ipld-prolly-indexer/) library

`cat example.csv | csv-to-prolly-db -o example.car`

```
NAME:
   csv-to-prolly-db - Ingest a CSV file into a prolly tree

USAGE:
   csv-to-prolly-db [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value      Path to save database CAR file to. (default: "./db.car")
   --input value, -i value       Path to CSV file to read. Omit to load from STDIN
   --collection value, -c value  Name of database collection to save into (default: "default")
   --help, -h                    show help
```
