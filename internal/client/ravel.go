package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cerebrotech/terraform-provider-ravel/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RavelClient struct {
	httpClient *resty.Client
}

func New(httpClient *resty.Client) *RavelClient {
	return &RavelClient{
		httpClient: httpClient,
	}
}

func (rc *RavelClient) CreateConfig(c context.Context, meta models.RavelConfigMeta, configFormat *models.RavelSchemaMeta, configDef map[string]any) (*models.RavelConfig, error) {
	ravelConfig := models.RavelConfig{
		Meta: meta,
		Spec: models.RavelConfigSpec{
			ConfigurationFormat: configFormat,
			Def:                 configDef,
		},
	}

	var definition []byte
	definition, err := json.Marshal(ravelConfig)
	if err != nil {
		return nil, err
	}
	tflog.Info(c, fmt.Sprintf("Create config: %s", string(definition)))
	res, err := rc.httpClient.R().SetContext(c).SetBody(ravelConfig).Post("/configurations")
	if err != nil {
		return nil, err
	}

	return rc.configProcess(res, err)
}

func (rc *RavelClient) DeleteConfig(c context.Context, configId string) error {
	res, err := rc.httpClient.R().SetContext(c).SetPathParams(map[string]string{
		"configId": configId,
	}).Delete("/configurations/{configId}")

	return rc.handleError(res, err)
}

func (rc *RavelClient) GetConfigVersion(c context.Context, configId string, version int) (*models.RavelConfig, error) {
	res, err := rc.httpClient.R().SetContext(c).SetPathParams(map[string]string{
		"configId": configId,
		"version":  strconv.Itoa(version),
	}).Get("/configurations/{configId}/versions/{version}?secrets=resolve")

	return rc.configProcess(res, err)
}

func (rc *RavelClient) configProcess(res *resty.Response, err error) (*models.RavelConfig, error) {
	if err := rc.handleError(res, err); err != nil {
		return nil, err
	}

	var config *models.RavelConfig
	if err := json.Unmarshal(res.Body(), &config); err != nil {
		return nil, err
	}

	return config, err
}

func (rc *RavelClient) handleError(res *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error communicating with Ravel. URL: %s - %d. Response: %s", res.Request.URL, res.StatusCode(), string(res.Body()))
	}

	return nil
}
