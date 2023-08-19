resource "nomad_job" "app" {
  jobspec = templatefile("${path.module}/guess-my-word.nomad", {
    image_name = var.artifact_url
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

variable "artifact_url" {
  description = "The artifact URL to use for the job"
  type        = string
}
