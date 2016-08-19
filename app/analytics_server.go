package app

import (
	"encoding/json"
	"net/http"
	"splash/processing/aggregation"
	"splash/services"
	"splash/logger"
)

type AnalyticsServer struct {
	config map[string]string
	aggregator *aggregation.Aggregator
	logger *logger.Logger
}

func NewAnalyticsServer(config map[string]string, aggregator *aggregation.Aggregator, logger *logger.Logger) *AnalyticsServer {
	return &AnalyticsServer{
		config: config,
		aggregator: aggregator,
		logger: logger,
	}
}

func (self *AnalyticsServer) Start() {

	listenAddr := self.config["host"] + self.config["port"]

	http.HandleFunc(self.config["path"], self.handler())

	self.logger.Alert("Analytics Server: ", http.ListenAndServe(listenAddr, nil))
}

func (self *AnalyticsServer) handler() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		self.respond(w)
	}
}

func (*AnalyticsServer) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Locking on Aggregated data
	dailyActiveUsers := <-aggregation.AggregationQueue

	serviceLocator := services.Locator{}
	logger := serviceLocator.Logger()

	logger.Info("finally", dailyActiveUsers)

	json.NewEncoder(w).Encode(map[string]interface{}{"dailyActiveUsers": dailyActiveUsers})
}
