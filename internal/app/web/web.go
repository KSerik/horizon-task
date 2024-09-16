package web

import (
	"encoding/json"
	"github.com/kserik/horizon-task/internal/pkg/repos"
	"log"
	"net/http"
	"strings"
)

func StartServer(addr string, dbClient repos.MarketplaceData) {
	http.HandleFunc("/aggregated-data", func(w http.ResponseWriter, r *http.Request) {
		data, err := dbClient.GetAggregatedData()
		if err != nil {
			http.Error(w, "Failed to get aggregated data", http.StatusInternalServerError)
			return
		}

		// removing time part from date
		for i := range data {
			data[i].Day = strings.TrimSuffix(data[i].Day, "T00:00:00Z")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}
