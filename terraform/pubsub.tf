# order service will send messages to this topic
locals {
  pubsub_service_account_roles = ["pubsub.publisher", "pubsub.subscriber"]
}

resource "google_pubsub_topic" "order_topic" {
  name    = "order-topic"
  project = google_project.go_pubsub.project_id
  labels = {
    app       = "go-pubsub"
    component = "order"
  }
  message_retention_duration = "86600s"
}

# subscriptions
resource "google_pubsub_subscription" "package_sub" {
  name    = "package_sub"
  project = google_project.go_pubsub.project_id
  topic   = google_pubsub_topic.order_topic.name

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages      = true
  ack_deadline_seconds       = 20

  expiration_policy {
    ttl = "300000.5s"
  }
  retry_policy {
    minimum_backoff = "10s"
  }
  enable_message_ordering = false

  labels = {
    app       = "go-pubsub"
    component = "packaging"
  }
}


resource "google_pubsub_subscription" "notification_sub" {
  name    = "notification_sub"
  project = google_project.go_pubsub.project_id
  topic   = google_pubsub_topic.order_topic.name

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages      = true
  ack_deadline_seconds       = 20

  expiration_policy {
    ttl = "300000.5s"
  }
  retry_policy {
    minimum_backoff = "10s"
  }
  enable_message_ordering = false

  labels = {
    app       = "go-pubsub"
    component = "notification"
  }
}

#SA
resource "google_service_account" "pubsub_service_account" {
  project      = google_project.go_pubsub.project_id
  account_id   = "pubsub-system"
  display_name = "pubsub-system"
}

resource "google_project_iam_member" "pubsub_project_sa_iam" {
  project = google_project.go_pubsub.project_id
  count   = length(local.pubsub_service_account_roles)
  role    = "roles/${local.pubsub_service_account_roles[count.index]}"
  member  = "serviceAccount:${google_service_account.pubsub_service_account.email}"
}
