# Sequence Take-Home Exercise: Blockchain Data Aggregator for Marketplace Analytics

## Description

This project is a data aggregation pipeline for marketplace analytics using blockchain data. It reads transaction data from a local or remote CSV file, processes, transforms, and stores the results in a ClickHouse database. The solution includes functionalities for extracting, normalizing, and calculating daily marketplace volume, daily transactions, aggregated volume data, and an API endpoint to get aggregated data for visualization.

## Solution Overview

I've implemented the solution in Go and Clickhouse. Data is processed in the Go program, but the aggregation is done in clickhouse using aggregation materialized views. I chose this method because aggregation in Clickhouse is faster and more efficient than Go. And we can later work with normalized and aggregated data in clickhouse. \
The Coingecko API fetches exchange rates for different currencies and converts the transaction values to USD. However, the coingecko API has a rate limit of 5 queries per minute. I've used an in-memory cache to store the exchange rates and added historical rates to accelerate the processing. To fetch the historical exchange prices from Coingecko API, you must provide coin IDs, which we can get from coin symbols. But the coin symbols are not always the same as the coin IDs, so I've used a hardcoded map in code to get the coin IDs.

For the project structure, I've used the standard project layout. The project is divided into five packages:
- bin - contains the compiled binary
- cmd - commands to run the application
- internal/app - packages for the applications
- internal/pkg - packages for shared code between applications 
- testdata - contains sample data

The program is a CLI application with four commands:
- migrate -  to create the tables in ClickHouse
- aggregate -  to read the data from a CSV file or URL, normalize it, and store it in ClickHouse.
- truncateall -  to truncate all data in the aggregates tables.
- web - to start an API server to fetch aggregated data for visualization

 The data is stored in ClickHouse in 3 tables: 
 - marketplace_data -  to store extracted and normalized transaction data. We may choose not to store this data, as the data is aggregated on insert.
 - agg_marketplace_data - to store pre-aggregated data. 
 - mv_agg_marketplace_data(materialized view) - materialized view aggregates data inserted in the marketplace_data table.

## Setup

### Prerequisites

- Go 1.16 or later
- Docker

### Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/kserik/horizon-task.git
    cd horizon-task
    ```

2. **Install dependencies**:
    ```sh
    go mod download
    ```

3. **Run ClickHouse in Docker**:
    ```sh
    docker run -d -p 18123:8123 -p 19000:9000 --name some-clickhouse-server --ulimit nofile=262144:262144 clickhouse/clickhouse-server
    ```
4. go build -o ./bin/horizon-task

### Configuration

- **ClickHouse Connection String**: Provide the ClickHouse connection string using the `--dsn` flag or `-d` flag. E.g. in docker: -d="clickhouse://default:@localhost:19000/default?compress=lz4"
- **File Path or URL**: Provide the local file path or URL to read data from using the `--filePath` or `-p` flag and `--fileURL` or `-u` flag respectively. E.g. -p="path/to/data.csv" or -u="http://example.com/data.csv"

### Usage
1. **Migrate the table schemas**:
    ```sh
    go run main.go migrate --dsn "clickhouse://default:@localhost:19000/default?compress=lz4"
   
    #or with binary
    ./bin/horizon-task migrate --dsn "clickhouse://default:@localhost:19000/default?compress=lz4"
   ```

2. **Run Data Aggregation**:
    ```sh
    go run main.go aggregate --dsn "clickhouse://default:@localhost:19000/default?compress=lz4" --filePath "path/to/data.csv"
    
    #Or using a URL:
    go run main.go aggregate --dsn "clickhouse://default:@localhost:19000/default?compress=lz4" --fileURL "http://example.com/data.csv"
    ```
3. **Start the API server**:
    ```sh
    go run main.go web --dsn "clickhouse://default:@localhost:19000/default?compress=lz4"
   
   #or with binary
   ./bin/horizon-task web --dsn "clickhouse://default:@localhost:19000/default?compress=lz4"
   ```

4. ** If needed truncate All Data**:
    ```sh
    go run main.go truncateall --dsn "clickhouse://default:@localhost:19000/default?compress=lz4"
    ```
### Notes

- Ensure that the ClickHouse server is running and accessible.
- The CSV file should have the following columns: `ts`, `event`, `project_id`, "currencySymbol" in `props` column , "currencyValueDecimal" in `vals` column.
- CurrencyValueDecimal of MATIC transactions doesn't look like a correct, numbers are too big
- I didn't implement timeouts for processing files in the application. If the data is too large, the application may hang. In production, we should implement timeouts using context and more error handling.
- I didn't imlement the deduplication of data, so if you run the aggregation command multiple times, the data will be duplicated. But this can be fixed by adding a unique id of transaction and using ReplacingMergeTree engine in ClickHouse.
