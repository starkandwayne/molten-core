package bucc

import (
	"strings"
)

const (
	rcTmpl = `
{
  "addons": [
    {
      "include": {
        "stemcell": [
          {
            "os": "ubuntu-trusty"
          },
          {
            "os": "ubuntu-xenial"
          }
        ]
      },
      "jobs": [
        {
          "name": "bosh-dns",
          "properties": {
            "api": {
              "client": {
                "tls": "((/dns_api_client_tls))"
              },
              "server": {
                "tls": "((/dns_api_server_tls))"
              }
            },
            "cache": {
              "enabled": true
            },
            "recursors": [
              "8.8.8.8",
              "8.8.4.4"
            ],
            "health": {
              "client": {
                "tls": "((/dns_healthcheck_client_tls))"
              },
              "enabled": true,
              "server": {
                "tls": "((/dns_healthcheck_server_tls))"
              }
            }
          },
          "release": "bosh-dns"
        }
      ],
      "name": "bosh-dns"
    },
    {
      "include": {
        "stemcell": [
          {
            "os": "windows2012R2"
          },
          {
            "os": "windows2016"
          },
          {
            "os": "windows1803"
          },
          {
            "os": "windows2019"
          }
        ]
      },
      "jobs": [
        {
          "name": "bosh-dns-windows",
          "properties": {
            "api": {
              "client": {
                "tls": "((/dns_api_client_tls))"
              },
              "server": {
                "tls": "((/dns_api_server_tls))"
              }
            },
            "cache": {
              "enabled": true
            },
            "health": {
              "client": {
                "tls": "((/dns_healthcheck_client_tls))"
              },
              "enabled": true,
              "server": {
                "tls": "((/dns_healthcheck_server_tls))"
              }
            }
          },
          "release": "bosh-dns"
        }
      ],
      "name": "bosh-dns-windows"
    }
  ],
  "releases": [
    {
      "name": "bosh-dns",
      "sha1": "d514ab3ae376778e106e17c22b78a8705690ae1d",
      "url": "https://bosh.io/d/github.com/cloudfoundry/bosh-dns-release?v=1.17.0",
      "version": "1.17.0"
    }
  ],
  "variables": [
    {
      "name": "/dns_healthcheck_tls_ca",
      "options": {
        "common_name": "dns-healthcheck-tls-ca",
        "is_ca": true
      },
      "type": "certificate"
    },
    {
      "name": "/dns_healthcheck_server_tls",
      "options": {
        "ca": "/dns_healthcheck_tls_ca",
        "common_name": "health.bosh-dns",
        "extended_key_usage": [
          "server_auth"
        ]
      },
      "type": "certificate"
    },
    {
      "name": "/dns_healthcheck_client_tls",
      "options": {
        "ca": "/dns_healthcheck_tls_ca",
        "common_name": "health.bosh-dns",
        "extended_key_usage": [
          "client_auth"
        ]
      },
      "type": "certificate"
    },
    {
      "name": "/dns_api_tls_ca",
      "options": {
        "common_name": "dns-api-tls-ca",
        "is_ca": true
      },
      "type": "certificate"
    },
    {
      "name": "/dns_api_server_tls",
      "options": {
        "ca": "/dns_api_tls_ca",
        "common_name": "api.bosh-dns",
        "extended_key_usage": [
          "server_auth"
        ]
      },
      "type": "certificate"
    },
    {
      "name": "/dns_api_client_tls",
      "options": {
        "ca": "/dns_api_tls_ca",
        "common_name": "api.bosh-dns",
        "extended_key_usage": [
          "client_auth"
        ]
      },
      "type": "certificate"
    }
  ]
}
`
)

func renderRuntimeConfig() string {
	raw := strings.ReplaceAll(rcTmpl, "\n", "")
	return raw
}
