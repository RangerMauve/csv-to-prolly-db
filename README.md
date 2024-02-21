# csv-to-prolly-db

Take a csv file and turn it into a prolly tree based database indexed by row number.

Based on the [ipld-prolly-indexer](https://github.com/RangerMauve/ipld-prolly-indexer/) library

Individual fields in a row will be parsed according to the [dag-json](https://ipld.io/docs/codecs/known/dag-json/) spec enabling structured data and lower memory usage in the case of numbers.

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

## nix flake

You can run and package `csv-to-prolly-db` with [nix](https://nixos.org/download/) by referencing this flake.

Run:

```bash
nix run github:RangerMauve/csv-to-prolly-db -- -i <my-csv>
```

Include as a package in another flake:

```nix
{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    csv-to-prolly-db.url = "github:RangerMauve/csv-to-prolly-db?ref=default";
  };

  outputs = inputs @ { self, ... }:
    (inputs.flake-utils.lib.eachDefaultSystem (system:
      let

        pkgs = import inputs.nixpkgs { inherit system; };

        csv-to-prolly-db = inputs.csv-to-prolly-db.packages.${system}.default;

    ...

}
```

