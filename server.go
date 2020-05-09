package main

import (
	"github.com/gorilla/mux"
	"log"
	"metricsserver/datastore"
	"metricsserver/handlers"
	"metricsserver/publish"
	"net/http"
)

func main() {
	sqlStore := datastore.NewSqlStore(10, 10)
	dataStore := datastore.NewLocalCache(sqlStore)
	publisher := publish.NewKafkaPublisher("localhost:9092", "metrics")
	metricHandler := handlers.NewMetricHandler(dataStore, publisher, 1000000, 250)
	router := mux.NewRouter()
	router.HandleFunc("/metrics", metricHandler.ProcessMetrics)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Print("Unable to start a http server.")
	}

}
