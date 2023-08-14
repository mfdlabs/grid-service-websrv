#  Copyright 2022 Nikita Petko <petko@vmminfra.net>
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

job "grid-service-websrv" {
  datacenters = ["*"]

  vault {
    policies = ["vault_secret_grid_service_websrv"]
  }

  # RBX Proxy nodes
  group "grid-service-websrv" {
    count = 5

    task "server" {
      driver = "docker"

      env {
        BIND_ADDRESS_IPV4  = "0.0.0.0"
        ENABLE_TLS_SERVER  = "false"
        INSECURE_PORT      = "${NOMAD_PORT_http}"
        DISABLE_IPV6       = "true"
        MFDLABS_ARC_SERVER = "true"
      }

      config {
        image = "mfdlabs/grid-service-websrv:latest"
      }

      resources {
        network {
          port "http" {}
        }
      }

      template {
        data = <<EOF
BIND_ADDRESS_IPv4 = ":{{ env "NOMAD_PORT_http" }}"
{{ with secret "kv-migration/grid-service-websrv" }}

{{ if .Data.data }}
{{ range $key, $value := .Data.data }}
{{ $key }} = "{{ $value }}"
{{ end }}
{{ end }}

{{ end }}
EOF

      service {
        name = "rbx-proxy"
        port = "http"

        tags = [
          "traefik.enable=true",
          "traefik.tags=http",
          "traefik.frontend.entryPoints=http,https",
          "traefik.frontend.passHostHeader=true",
          "traefik.frontend.rule=HostRegexp:(avatar|clientsettingscdn).sitetest4.robloxlabs.com,(ephemeralcounters|versioncompatibility).api.sitetest4.robloxlabs.com",
          "traefik.backend.loadbalancer.method=wrr",
          "traefik.backend.buffering.retryExpression=IsNetworkError() && Attempts() <= 2"
        ]

        check {
          type     = "http"
          path     = "/_lb/_/health"
          interval = "2s"
          timeout  = "2s"
        }
      }
    }
  }
}