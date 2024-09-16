-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS agg_marketplace_data
(
    day              Date,
    project_id       Int64,
    no_of_trx        AggregateFunction(count, UInt64),
    total_volume_usd AggregateFunction(sum, Decimal(39, 18))
) ENGINE = AggregatingMergeTree
ORDER BY (day, project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS agg_marketplace_data;
-- +goose StatementEnd
