locals {
  services = ["pubsub.googleapis.com"]
}

resource "google_project" "go_pubsub" {
  name            = "prj-go-pubsub"
  project_id      = "prj-go-pubsub"
  billing_account = var.billing_account
}

# enable apis
resource "google_project_service" "go_pubsub_services" {
  count                      = length(local.services)
  project                    = google_project.go_pubsub.project_id
  service                    = local.services[count.index]
  disable_dependent_services = true
}

resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name          = "${random_id.bucket_prefix.hex}-bucket-tfstate"
  project       = google_project.go_pubsub.project_id
  force_destroy = true # this is just a demo want it to be easy to clean up
  location      = "AUSTRALIA-SOUTHEAST2"
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}