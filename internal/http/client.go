package http

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

var GlobalHTTPClient = &http.Client{}

func New(url, version, token string) *resty.Client {
	client := resty.NewWithClient(GlobalHTTPClient)
	client.SetBaseURL(url)
	client.SetHeader("User-Agent", fmt.Sprintf("domino/terraform-provider-ravel:%s", version))
	client.SetHeader("X-Request-Id", uuid.NewString())
	client.SetHeader("X-Api-Token", token)

	return client
}
