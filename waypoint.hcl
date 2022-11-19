project = "guess-my-word"

app "web" {
  build {
    use "docker-pull" {
      disable_entrypoint = true
      image              = "ghcr.io/taiidani/guess-my-word"
      tag                = "latest"
    }
  }

  deploy {
    use "nomad-jobspec" {
      jobspec = templatefile("${path.app}/.github/guess-my-word.hcl", {
        artifact = {
          image = "ghcr.io/taiidani/guess-my-word"
          tag   = "latest"
        },
      })
    }
  }
}
