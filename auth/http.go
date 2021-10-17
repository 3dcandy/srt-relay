package auth

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/voc/srtrelay/stream"
)

type httpAuth struct {
	config HTTPAuthConfig
	client *http.Client
}

type HTTPAuthConfig struct {
	URL           string
	Application   string
	Timeout       time.Duration // Timeout for Auth request
	PasswordParam string        // POST Parameter containing stream passphrase
}

// NewHttpAuth creates an Authenticator with a HTTP backend
func NewHTTPAuth(config HTTPAuthConfig) *httpAuth {
	return &httpAuth{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Implement Authenticator

// Authenticate sends form-data in a POST-request to the configured url.
// If the response code is 2xx the publish/play is allowed, otherwise it is denied.
// This should be compatible with nginx-rtmps on_play/on_publish directives.
// https://github.com/arut/nginx-rtmp-module/wiki/Directives#on_play
func (h *httpAuth) Authenticate(streamid stream.StreamID) bool {
	response, err := h.client.PostForm(h.config.URL, url.Values{
		"call":                 {streamid.Mode().String()},
		"app":                  {h.config.Application},
		"name":                 {streamid.Name()},
		h.config.PasswordParam: {streamid.Password()},
	})

	if err != nil {
		log.Println("http-auth:", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return false
	}

	return true
}
