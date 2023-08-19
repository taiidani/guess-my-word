variable "artifact_url" {
  type = string
}

job "guess-my-word" {
  datacenters = ["dc1"]
  type        = "service"

  update {
    min_healthy_time  = "30s"
    healthy_deadline  = "1m"
    progress_deadline = "10m"
    auto_revert       = false
  }

  group "app" {
    task "app" {
      driver = "docker"

      config {
        image = var.artifact_url
        args  = ["/app"]
        ports = ["web"]
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

        // check_restart {
        //   limit           = 3
        //   grace           = "15s"
        //   ignore_warnings = false
        // }
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

    task "redirect" {
      driver = "docker"

      config {
        image = "containous/whoami"
        args  = ["-port=81"]
        ports = ["redirect"]
      }

      service {
        name     = "guess-my-word-redirect"
        provider = "nomad"
        port     = "redirect"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.guessredirect.rule=Host(`guess.taiidani.com`)",
          "traefik.http.routers.guessredirect.middlewares=guessredirect@nomad",
          "traefik.http.routers.guessredirectsecure.rule=Host(`guess.taiidani.com`)",
          "traefik.http.routers.guessredirectsecure.tls=true",
          "traefik.http.routers.guessredirectsecure.tls.certresolver=le",
          "traefik.http.routers.guessredirectsecure.middlewares=guessredirect@nomad",
          "traefik.http.middlewares.guessredirect.redirectregex.regex=^http.?://guess.taiidani.com/(.*)",
          "traefik.http.middlewares.guessredirect.redirectregex.replacement=https://guessmyword.xyz/",
        ]

        // check_restart {
        //   limit           = 3
        //   grace           = "15s"
        //   ignore_warnings = false
        // }
      }

      resources {
        cpu    = 25
        memory = 32
      }
    }

    network {
      mode = "bridge"
      port "web" { to = 3000 }
      port "redirect" { to = 81 }
    }

    vault {
      policies = ["hcp-root"]
    }
  }
}
