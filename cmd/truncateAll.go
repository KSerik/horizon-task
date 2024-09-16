/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/kserik/horizon-task/internal/pkg/repos/clickhouse"
	"log"

	"github.com/spf13/cobra"
)

// truncateAllCmd represents the truncateAll command
var truncateAllCmd = &cobra.Command{
	Use:   "truncateall",
	Short: "Truncate all data in aggregates tables",
	Run: func(cmd *cobra.Command, args []string) {
		if dsn == "" {
			log.Fatalln("Please provide a Clickhouse connection string")
		}
		chConn, err := clickhouse.NewClickhouseConn(dsn)
		if err != nil {
			log.Fatalf("Couldn't create clickhouse client: %s\n", err)
		}
		defer chConn.Close()

		repo := clickhouse.NewMarketplaceDataRepo(chConn)

		if err := repo.TruncateAllData(); err != nil {
			log.Fatalln("Couldn't truncate data: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(truncateAllCmd)
}
