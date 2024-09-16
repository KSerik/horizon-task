package cmd

import (
	"github.com/kserik/horizon-task/internal/app/migrate"
	"github.com/kserik/horizon-task/internal/pkg/repos/clickhouse"
	"log"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate all aggregation tables in a new DB",
	Run: func(cmd *cobra.Command, args []string) {
		if dsn == "" {
			log.Fatalln("Please provide a Clickhouse connection string")
		}

		chClient, err := clickhouse.NewClickhouseConn(dsn)
		if err != nil {
			log.Fatalf("Couldn't create clickhouse client: %s\n", err)
		}
		defer chClient.Close()

		app := migrate.NewApp(chClient)

		if err := app.Up(); err != nil {
			log.Fatalln("Couldn't migrating tables:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
