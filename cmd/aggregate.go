package cmd

import (
	"encoding/csv"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kserik/horizon-task/internal/app/aggregate"
	"github.com/kserik/horizon-task/internal/pkg/exchange"
	"github.com/kserik/horizon-task/internal/pkg/repos/clickhouse"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
)

var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Run the data aggregation pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" && fileURL == "" {
			log.Fatalln("Please provide either a filePath or a URL")
		}

		if dsn == "" {
			log.Fatalln("Please provide a Clickhouse connection string")
		}

		chConn, err := clickhouse.NewClickhouseConn(dsn)
		if err != nil {
			log.Fatalf("Coudln't create clickhouse client: %s\n", err)
		}
		defer chConn.Close()

		repo := clickhouse.NewMarketplaceDataRepo(chConn)
		app := aggregate.NewApp(exchange.NewService(), repo)

		// if local file provided it will be used, otherwise file will be downloaded from URL
		fileReader, err := prepareFileReader(filePath, fileURL)
		if err != nil {
			log.Fatalf("Error preparing file reader: %s\n", err)
		}
		defer fileReader.Close()

		csvReader := csv.NewReader(fileReader)

		if err := app.Run(csvReader); err != nil {
			log.Fatalln("Couldn't run the aggregation pipeline: ", err)
		}
	},
}

func prepareFileReader(filePath, urlFilePath string) (io.ReadCloser, error) {
	var err error
	// Open the CSV filePath
	var file io.ReadCloser
	if filePath != "" {
		file, err = os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening filePath: %w", err)
		}
	} else {
		file, err = getFileFromURL(urlFilePath)
		if err != nil {
			return nil, fmt.Errorf("error getting file from URL: %w", err)
		}
	}

	return file, nil
}

func getFileFromURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("making GET request: %w", err)
	}

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status is not OK: %s", resp.Status)
	}

	return resp.Body, nil
}

var (
	filePath string
	fileURL  string
)

func init() {
	rootCmd.AddCommand(aggregateCmd)
	aggregateCmd.Flags().StringVarP(&filePath, "filePath", "p", "", "Local file path to read data from")
	aggregateCmd.Flags().StringVarP(&fileURL, "fileURL", "u", "", "file URL to read data from")
}
