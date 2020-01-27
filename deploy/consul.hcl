service {
    name = "guess-my-word"
    port = 3000
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