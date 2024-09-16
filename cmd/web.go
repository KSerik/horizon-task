package cmd

import (
	"github.com/kserik/horizon-task/internal/app/web"
	"github.com/kserik/horizon-task/internal/pkg/repos/clickhouse"
	"log"

	"github.com/spf13/cobra"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Starts the web server",
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

		web.StartServer(":8080", repo)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
