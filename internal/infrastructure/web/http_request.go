package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/felisest/comproxy/internal/infrastructure/config"
)

type HttpRequest struct {
	cfg config.Config
}

func NewHttpRequest(cfg *config.Config) *HttpRequest {
	return &HttpRequest{
		cfg: *cfg,
	}
}

func (r *HttpRequest) Post(request []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", r.cfg.TestingHost+r.cfg.Path, bytes.NewBuffer(request))
	if err != nil {

		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
