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
    task "web" {
      driver = "docker"

      config {
        image = "${artifact.image}:${artifact.tag}"
        ports = ["web"]
      }

      service {
        name     = "guess-my-word-web"
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
        cpu    = 25
        memory = 32
      }
    }

    task "api" {
      driver = "docker"

      config {
        image = "${artifact.image}:${artifact.tag}"
        args  = ["/app"]
        ports = ["api"]
      }

      env {
        ADDR       = "0.0.0.0"
        GIN_MODE   = "release"
        REDIS_ADDR = "127.0.0.1:$${NOMAD_PORT_redis}"
        ORIGIN     = "guessmyword.xyz"
      }

      service {
        name     = "guess-my-word-api"
        provider = "nomad"
        port     = "api"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.guesssecureapi.rule=Host(`guessmyword.xyz`) && PathPrefix(`/api`)",
          "traefik.http.routers.guesssecureapi.tls=true",
          "traefik.http.routers.guesssecureapi.tls.certresolver=le",
          "traefik.http.routers.guesssecureapi.middlewares=guess@nomad",
        ]

        // check_restart {
        //   limit           = 3
        //   grace           = "15s"
        //   ignore_warnings = false
        // }
      }

      resources {
        cpu    = 50
        memory = 64
      }
    }

    volume "hashistack" {
      type      = "host"
      source    = "hashistack"
      read_only = "false"
    }

    task "redis" {
      driver = "docker"

      config {
        image = "redis:6"
        args  = ["redis-server", "--dir", "/data/guess"]
        ports = ["redis"]
      }

      volume_mount {
        volume      = "hashistack"
        destination = "/data"
        read_only   = "false"
      }

      resources {
        cpu    = 25
        memory = 32
      }
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
      port "web" { to = 80 }
      port "api" { to = 3000 }
      port "redirect" { to = 81 }
      port "redis" { to = 6379 }
    }
  }
}
