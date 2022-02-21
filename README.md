# tractor-beam 🛸

An application for synchronising key-value data from a spreadsheet into a Redis database.

Given a spreadsheet with two headered columns representing a key and a value, tractor-beam will save the key-value pairs to Redis.

## Usage

Install the application using the `go` CLI or download the appropriate binary from the [releases page.](https://github.com/djaustin/tractor-beam/releases)

```bash
go install github.com/djaustin/tractor-beam@latest
```

### One-time Sync

tractor-beam can be used to complete a single synchronisation of the database from a file. 

```bash
tractor-beam sync <spreadsheet-file> <redis-address>
```

```bash
tractor-beam sync Book.xlsx localhost:6379
```

### Watch Mode

Given a spreadsheet file, tractor-beam can continually monitor the file for changes and keep the Redis instance up to date when any changes occur.

```bash
tractor-beam watch <spreadsheet-file> <redis-address>
```

```bash
tractor-beam watch Book.xlsx localhost:6379
```

## Configuration

The tractor-beam application can be configured using the following methods ranked in order of precedence (decreasing)

1. Command line flags
2. Environment variables
3. Configuration file

### Command line flags

In order to view the possible command line flags for a command, use the `-h` flag to view the help message.

```bash
tractor-beam -h
tractor-beam sync -h
tractor-beam watch -h
```

### Configuration Files

By default, tractor-beam looks for a file named either `tractor-beam` or `tractor-beam.yaml` in the following locations:
* The directory from which tractor-beam is run
* The user's home directory
* /etc/tractor-beam/

An reference configuration file is included in this repository ([tractor-beam.reference.yaml](tractor-beam.reference.yaml))
