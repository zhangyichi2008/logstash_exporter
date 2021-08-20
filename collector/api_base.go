package collector

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
)

var (
	promlogCfg = &promlog.Config{}
	logger     = promlog.New(promlogCfg)
)

// HTTPHandler type
type HTTPHandler struct {
	Endpoint string
}

// Get method for HTTPHandler
func (h *HTTPHandler) Get() (http.Response, error) {
	response, err := http.Get(h.Endpoint)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}

// HTTPHandlerInterface interface
type HTTPHandlerInterface interface {
	Get() (http.Response, error)
}

func getMetrics(h HTTPHandlerInterface, target interface{}) error {
	response, err := h.Get()
	if err != nil {
		level.Error(logger).Log("Cannot retrieve metrics: %v", err)
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			level.Error(logger).Log("Cannot close response body: %v", err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(target)
	if err != nil {
		level.Error(logger).Log("Cannot parse Logstash response json: %v", err)
	}

	return err
}
