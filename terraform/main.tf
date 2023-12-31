resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name          = "${random_id.bucket_prefix.hex}-bucket-tfstate"
  force_destroy = true # this is just a demo want it to be easy to clean up
  location      = "AUSTRALIA-SOUTHEAST2"
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}