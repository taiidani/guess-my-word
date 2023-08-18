resource "nomad_job" "app" {
  jobspec = templatefile("${path.module}/guess-my-word.nomad", {
    version = var.version
  })
  detach = false

  hcl2 {
    enabled = true
  }
}

terraform {
  required_version = ">= 1.4.0"

  required_providers {
    nomad = {
      source  = "hashicorp/nomad"
      version = "1.4.20"
    }
  }

  cloud {
    organization = "rnd"

    workspaces {
      name = "guess-my-word"
    }
  }
}

provider "nomad" {
  address = "http://127.0.0.1:4646"
}

variable "version" {
  description = "The artifact version to use for the job"
  type        = string
}
