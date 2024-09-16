package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/kserik/horizon-task/internal/pkg/repos"
	"time"
)

type MarketplaceDataRepo struct {
	conn *sql.DB
}

func NewClickhouseConn(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	return conn, conn.Ping()
}

func NewMarketplaceDataRepo(conn *sql.DB) repos.MarketplaceData {
	return &MarketplaceDataRepo{conn: conn}
}

func (m *MarketplaceDataRepo) InsertTransactions(transactions []repos.Transaction) error {
	scope, err := m.conn.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer scope.Rollback()

	batch, err := scope.Prepare("INSERT INTO marketplace_data (ts, event, project_id, currency_symbol, currency_value_decimal, usd_conversion_rate, usd_amount) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer batch.Close()

	for _, trx := range transactions {
		if _, err := batch.Exec(trx.Timestamp, trx.Event, trx.ProjectID, trx.CurrencySymbol, trx.CurrencyValue, trx.USDConversionRate, trx.USDValue); err != nil {
			return fmt.Errorf("error executing batch: %w", err)
		}
	}

	if err := scope.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (m *MarketplaceDataRepo) TruncateAllData() error {
	if _, err := m.conn.Exec("TRUNCATE TABLE marketplace_data"); err != nil {
		return err
	}

	if _, err := m.conn.Exec("TRUNCATE TABLE agg_marketplace_data"); err != nil {
		return err
	}

	return nil
}

func (m *MarketplaceDataRepo) GetAggregatedData() ([]repos.AggregatedData, error) {
	rows, err := m.conn.Query(`SELECT day,
       project_id,
       countMerge(no_of_trx)      as number_of_transactions,
       sumMerge(total_volume_usd) as total_volume_usd
FROM agg_marketplace_data
GROUP By day, project_id
ORDER By day, project_id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []repos.AggregatedData
	for rows.Next() {
		var agd repos.AggregatedData
		if err := rows.Scan(&agd.Day, &agd.ProjectID, &agd.NoOfTransactions, &agd.TotalUSDVolume); err != nil {
			return nil, err
		}
		data = append(data, agd)
	}

	return data, nil
}
