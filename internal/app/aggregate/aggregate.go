package aggregate

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kserik/horizon-task/internal/pkg/exchange"
	"github.com/kserik/horizon-task/internal/pkg/repos"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	ExchangeService *exchange.Service
	chClient        repos.MarketplaceData
}

func NewApp(exchangeService *exchange.Service, repo repos.MarketplaceData) *App {
	return &App{
		ExchangeService: exchangeService,
		chClient:        repo,
	}
}

func (a *App) Run(csvReader *csv.Reader) error {
	// Read the first record which contains the headersS
	_, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV filePath: %w", err)
	}

	var transactions []repos.Transaction
	for {
		var record []string
		record, err = csvReader.Read()
		if err != nil {
			break
		}

		transaction, err := a.parseRowToTransaction(record)
		if err != nil {
			return fmt.Errorf("error parsing record to transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading CSV records: %w", err)
	}

	if len(transactions) == 0 {
		return fmt.Errorf("no transactions to insert")
	}

	if err = a.chClient.InsertTransactions(transactions); err != nil {
		return fmt.Errorf("error inserting transactions: %w", err)
	}

	return nil
}

type propColumn struct {
	CurrencySymbol string `json:"currencySymbol"`
}

type numsColumn struct {
	CurrencyValueDecimal string `json:"currencyValueDecimal"`
}

func (a *App) parseRowToTransaction(record []string) (repos.Transaction, error) {
	var prop propColumn
	var nums numsColumn
	err := json.Unmarshal([]byte(record[14]), &prop)
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error unmarshalling props column: %w", err)
	}

	err = json.Unmarshal([]byte(record[15]), &nums)
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error unmarshalling nums column: %w", err)
	}

	projectID, err := strconv.Atoi(record[3])
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error converting project ID to int: %w", err)
	}
	currVal, err := decimal.NewFromString(nums.CurrencyValueDecimal)
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error converting currency value to decimal: %w", err)
	}

	ts, err := time.Parse("2006-01-02 15:04:05.000", record[1])
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error parsing timestamp: %w", err)
	}

	exchangeRate, err := a.ExchangeService.GetExchangeRate(prop.CurrencySymbol, ts.Format("02-01-2006"))
	if err != nil {
		return repos.Transaction{}, fmt.Errorf("error getting exchange rate: %w", err)
	}

	// we are doing the calculations here to control overflow precision, otherwise we can do it in clickhouse
	usdAmount := currVal.Mul(exchangeRate)

	return repos.Transaction{
		Timestamp:         ts,
		Event:             record[2],
		ProjectID:         projectID,
		CurrencySymbol:    prop.CurrencySymbol,
		CurrencyValue:     currVal,
		USDConversionRate: exchangeRate,
		USDValue:          usdAmount,
	}, nil
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
