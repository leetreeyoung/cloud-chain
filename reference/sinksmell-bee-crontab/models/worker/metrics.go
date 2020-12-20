package worker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitPromMetrics  初始化 prometheus
func InitPromMetrics() (err error) {

	if Conf.PromPort == 0 {
		return errors.New("prometheus port is empty")
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", Conf.PromPort), nil)
	}()

	return nil
}
