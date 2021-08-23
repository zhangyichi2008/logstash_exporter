package collector

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Cannot retrieve metrics.")
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Cannot close response body")
		}
	}()

	err = json.NewDecoder(response.Body).Decode(target)
	if err != nil {
		log.Error().Err(err).Msg("Cannot parse logstash response JSON")
	}

	return err
}
