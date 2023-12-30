terraform {
  backend "gcs" {
    # todo
    # bucket = "replaceme"
    prefix = "terraform/state"
  }
}