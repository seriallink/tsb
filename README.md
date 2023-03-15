# tsb (Timescale Benchmarking)

Cloud Engineer Assignment - Benchmarking

## Usage

    tsb --input data/query_params.csv --workers 10

### Flags

```
Usage:
  tsb [flags]

Flags:
  -h, --help           help for tsb
  -i, --input string   csv file name (leave empty to use stdin)
  -w, --workers int    number of concurrent workers (default 10)
```

## Steps to run the assignment

1. Build container using Dockerfile
```
docker build -t tsb_image .
```

2. Run container:
```
docker run -d --name tsb_container tsb_image
```

3. From container terminal, run migration:
```
./migration.sh
```

4. From container terminal, run the tool:
```
./tsb -i data/query_params.csv -w 10
```

5. Output will be displayed in the terminal like this:
```
number of queries processed: 200
total processing time across all queries: 928.786ms
minimum query time (for a single query): 3.7932ms
maximum query time: 6.3703ms
median query time: 3.60055ms
average query time: 4.64393ms
```
