package prometheusquerier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// QueryInstant performs an instant query against a Prometheus server.
func QueryTSDBHeadSeries(cluster string) (float64, error) {
	config := api.Config{
		Address:      fmt.Sprintf("http://prometheus-operated.%s-prometheus.svc:9090/%s", cluster, cluster),
		RoundTripper: http.DefaultTransport,
	}

	// Create new client.
	c, err := api.NewClient(config)
	if err != nil {
		return 0, err
	}

	// Run query against client.
	api := v1.NewAPI(c)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	val, _, err := api.Query(ctx, "max_over_time(prometheus_tsdb_head_series[6h])", time.Now()) // Ignoring warnings for now.
	cancel()
	if err != nil {
		return 0, err
	}

	switch val.Type() {
	case model.ValVector:
		vector := val.(model.Vector)
		return float64(vector[0].Value), nil
	default:
		return 0, errors.New("failed to get current number of time series")
	}
}
