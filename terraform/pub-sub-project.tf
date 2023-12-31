locals {
  services = ["pubsub.googleapis.com"]
}

resource "google_project" "go_pubsub" {
  name       = "go-pubsub"
  project_id = "go-pubsub"
}

# enable apis
resource "google_project_service" "go_pubsub_services" {
  count                      = length(local.services)
  project                    = google_project.go_pubsub.project_id
  service                    = local.services[count.index]
  disable_dependent_services = true
}