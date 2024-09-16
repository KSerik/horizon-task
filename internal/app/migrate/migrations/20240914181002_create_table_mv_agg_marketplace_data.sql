-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_agg_marketplace_data TO agg_marketplace_data
AS
(
SELECT toDate(ts)           AS day,
       project_id           as project_id,
       countState(*)        as no_of_trx,
       sumState(usd_amount) AS total_volume_usd
FROM marketplace_data
GROUP BY day, project_id
order by (day, project_id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mv_agg_marketplace_data;
-- +goose StatementEnd
