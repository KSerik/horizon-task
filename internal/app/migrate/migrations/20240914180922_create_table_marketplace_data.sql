-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS marketplace_data
(
    ts                     DATETIME64(3, 'UTC'),
    event                  LowCardinality(String),
    project_id             Int64,
    currency_symbol        LowCardinality(String),
    currency_value_decimal Decimal(39, 18),
    usd_conversion_rate    Decimal(39, 18),
    usd_amount             Decimal(39, 18)
) ENGINE = MergeTree()
      ORDER BY (ts, project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS marketplace_data;
-- +goose StatementEnd
