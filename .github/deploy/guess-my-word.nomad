variable "artifact_url" {
  type = string
}

job "guess-my-word" {
  datacenters = ["dc1"]
  type        = "service"

  update {
    healthy_deadline  = "1m"
    progress_deadline = "2m"
    auto_revert       = true
  }

  reschedule {
    attempts  = 0
    unlimited = false
  }

  group "app" {
    count = 2

    task "app" {
      driver = "exec"

      config {
        command = "guess-my-word"
        args    = ["--port=${NOMAD_PORT_web}"]
      }

      artifact {
        source = var.artifact_url
      }

      env {
        ADDR     = "0.0.0.0"
        GIN_MODE = "release"
        ORIGIN   = "guessmyword.xyz"
      }

      template {
        data        = <<EOF
            REDIS_HOST="{{with secret "credentials/digitalocean/redis"}}{{ .Data.data.private_host }}{{end}}"
            REDIS_PORT="{{with secret "credentials/digitalocean/redis"}}{{ .Data.data.port }}{{end}}"
            REDIS_USER="{{with secret "credentials/digitalocean/redis"}}{{ .Data.data.user }}{{end}}"
            REDIS_PASSWORD="{{with secret "credentials/digitalocean/redis"}}{{ .Data.data.password }}{{end}}"
            REDIS_DB=1
        EOF
        destination = "${NOMAD_SECRETS_DIR}/secrets.env"
        env         = true
      }

      service {
        name     = "guess-my-word"
        provider = "nomad"
        port     = "web"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.guess.rule=Host(`guessmyword.xyz`)",
          "traefik.http.routers.guess.middlewares=guess@nomad",
          "traefik.http.routers.guesssecure.rule=Host(`guessmyword.xyz`)",
          "traefik.http.routers.guesssecure.tls=true",
          "traefik.http.routers.guesssecure.tls.certresolver=le",
          "traefik.http.routers.guesssecure.middlewares=guess@nomad",
          "traefik.http.middlewares.guess.redirectscheme.permanent=true",
          "traefik.http.middlewares.guess.redirectscheme.scheme=https",
        ]

        check_restart {
          limit           = 3
          grace           = "15s"
          ignore_warnings = false
        }
      }

      resources {
        cpu    = 50
        memory = 128
      }
    }

    volume "hashistack" {
      type      = "host"
      source    = "hashistack"
      read_only = "false"
    }

    network {
      port "web" {}
    }

    vault {
      policies = ["hcp-root"]
    }
  }
}
