package resources_test

import (
	"github.com/cerebrotech/terraform-provider-ravel/internal/provider"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"net/http"
	"net/http/httptest"
	"testing"
)

const smtpConfigurationDirectory = "../../it/stmp_configuration"

var getAPIResponse = `{
    "id": "7eb918e0-49b6-4519-bb5a-850c42d8da04",
    "created_at": 1698258912,
    "meta": {
        "name": "domino-cloud-smtp-configuration-test",
        "scope": {
            "category": "fleetcommand-configuration-manager",
            "fleetcommand_account": "ldebello-account",
            "type": "configuration"
        },
        "version": 0,
        "labels": {
            "minSchemaVersion": "001.000"
        }
    },
    "spec": {
        "configurationFormat": {
            "name": "email",
            "scope": {
                "category": "fleetcommand-configuration-manager",
                "source": "domino/release",
                "type": "schema"
            },
            "version": "1.0.0"
        },
        "def": {
            "email_notifications": {
                "authentication": {
                    "password": "password",
                    "username": "user"
                },
                "enable_ssl": true,
                "enabled": true,
                "from_address": "cloud-support@dominodatalab.com",
                "port": 465,
                "server": "email-smtp.us-east-1.amazonaws.com"
            }
        }
    }
}`

var createAPIResponse = `{
    "id": "7eb918e0-49b6-4519-bb5a-850c42d8da04",
    "created_at": 1698258912,
    "meta": {
        "name": "domino-cloud-smtp-configuration-test",
        "scope": {
            "category": "fleetcommand-configuration-manager",
            "fleetcommand_account": "ldebello-account",
            "type": "configuration"
        },
        "version": 0,
        "labels": {
            "minSchemaVersion": "001.000"
        }
    },
    "spec": {
        "configurationFormat": {
            "name": "email",
            "scope": {
                "category": "fleetcommand-configuration-manager",
                "source": "domino/release",
                "type": "schema"
            },
            "version": "1.0.0"
        },
        "def": {
            "email_notifications": {
                "authentication": {
                    "password": "secret://default/7ed44ae0-2313-4e00-a286-be15c78724a4",
                    "username": "user"
                },
                "enable_ssl": true,
                "enabled": true,
                "from_address": "cloud-support@dominodatalab.com",
                "port": 465,
                "server": "email-smtp.us-east-1.amazonaws.com"
            }
        }
    }
}`

func mockAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if r.Method == "GET" {
		_, err := w.Write([]byte(getAPIResponse))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else if r.Method == "POST" {
		_, err := w.Write([]byte(createAPIResponse))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func TestNewConfigurationResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockAPIHandler))
	defer server.Close()

	t.Setenv(provider.RavelToken, uuid.New().String())
	t.Setenv(provider.RavelURL, server.URL)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"ravel": providerserver.NewProtocol6WithError(provider.New("development")()),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory(smtpConfigurationDirectory),
			},
		},
	})
}
