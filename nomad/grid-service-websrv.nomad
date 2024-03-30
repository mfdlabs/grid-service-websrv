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
    count = 3

    network {
      mode = "host"

      port "http" {}
    }

    task "server" {
      driver = "docker"

      config {
        image        = "mfdlabs/grid-service-websrv:latest"
        network_mode = "host"
        ports        = ["http"]
      }

      template {
        data        = <<EOF
BIND_ADDRESS_IPv4 = ":{{ env "NOMAD_PORT_http" }}"

{{ range service "vault" }}
VAULT_ADDR = "http://{{ .Address }}:8200"
{{ end }}

{{ with secret "kv-migration/grid-service-websrv" }}

{{ if .Data.data }}
{{ range $key, $value := .Data.data }}
{{ $key }} = "{{ $value }}"
{{ end }}
{{ end }}

{{ end }}
EOF
        destination = "secrets/grid-service-websrv.env"
        env         = true
      }

      service {
        name = "grid-service-websrv"
        port = "http"

        tags = [
          "traefik.enable=true",
          "traefik.tags=http",
          "traefik.frontend.entryPoints=http,https",
          "traefik.frontend.passHostHeader=true",
          "traefik.frontend.rule=HostRegexp:{host:(.+)}.sitetest4.robloxlabs.com,{host:(.+)}.api.sitetest4.robloxlabs.com",
          "traefik.backend.loadbalancer.method=wrr",
          "traefik.backend.buffering.retryExpression=IsNetworkError() && Attempts() <= 2"
        ]

        check {
          type     = "http"
          path     = "/metrics"
          interval = "2s"
          timeout  = "2s"
        }
      }
    }
  }
}
