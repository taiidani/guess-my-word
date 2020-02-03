job "guess-my-word" {
  datacenters = ["pi"]
  type        = "service"

    task "guess-my-word" {
        driver = "exec"

        config {
            command = "/usr/local/bin/guess-my-word"
        }

        env {
          ADDR = "0.0.0.0"
          GO_ENV = "production"
        }

        service {
            port = "http"
            tags = [
                "traefik.enable=true",
                "traefik.http.routers.guess.rule=Host(`guess.home.ryannixon.com`)",
                "traefik.http.routers.guess.middlewares=guess@consulcatalog",
                "traefik.http.routers.guesssecure.rule=Host(`guess.home.ryannixon.com`)",
                "traefik.http.routers.guesssecure.tls=true",
                "traefik.http.routers.guesssecure.tls.certresolver=le",
                "traefik.http.routers.guesssecure.middlewares=guess@consulcatalog",
                "traefik.http.middlewares.guess.redirectscheme.permanent=true",
                "traefik.http.middlewares.guess.redirectscheme.scheme=https",
            ]
        }

        resources {
          network {
            port "http" {
              static = 3000
            }
          }
        }
    }
}
