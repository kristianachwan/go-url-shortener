package server

import "github.com/prometheus/client_golang/prometheus"

var (
	apiCallsByRoute = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls_by_route",
			Help: "Number of API calls by route",
		},
		[]string{"route"},
	)
)

func init() {
	prometheus.MustRegister(apiCallsByRoute)

}
