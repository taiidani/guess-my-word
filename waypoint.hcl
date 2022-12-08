project = "guess-my-word"

app "web" {
  build {
    hook {
      when       = "before"
      command    = ["go", "build", "-o", "./guess-my-word"]
      on_failure = "fail"
    }

    use "docker" {
      disable_entrypoint = true
      target             = "dist"
    }

    registry {
      use "docker" {
        image = "ghcr.io/taiidani/guess-my-word"
        tag   = "latest"
      }
    }

    hook {
      when       = "after"
      command    = ["rm", "./guess-my-word"]
      on_failure = "continue"
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
