resource "random_id" "random_suffix" {
  byte_length = 8
}

data "google_project" "project" {}

data "archive_file" "archive" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "../tmp/function.zip"
}

resource "google_storage_bucket" "source" {
  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }
  location = var.location
  name     = "source-${random_id.random_suffix.hex}"
}

resource "google_storage_bucket_object" "source" {
  bucket = google_storage_bucket.source.name
  name   = data.archive_file.archive.output_sha256
  source = "../tmp/function.zip"
}

resource "google_service_account" "function_runner" {
  account_id = "function-runner"
}

resource "google_project_iam_member" "function_runner" {
  member  = "serviceAccount:${google_service_account.function_runner.email}"
  project = data.google_project.project.project_id
  role    = "roles/secretmanager.secretAccessor"
}

resource "google_cloudfunctions2_function" "function" {
  location = var.location
  name     = var.project
  build_config {
    runtime     = "go122"
    entry_point = "EntryPoint"
    source {
      storage_source {
        bucket = google_storage_bucket.source.name
        object = google_storage_bucket_object.source.name
      }
    }
  }
  service_config {
    secret_environment_variables {
      project_id = data.google_project.project.project_id
      version    = "latest"
      secret     = google_secret_manager_secret.mailgun_api_key.secret_id
      key        = "MAILGUN_API_KEY"
    }
    secret_environment_variables {
      project_id = data.google_project.project.project_id
      version    = "latest"
      secret     = google_secret_manager_secret.mailgun_domain.secret_id
      key        = "MAILGUN_DOMAIN"
    }
    secret_environment_variables {
      project_id = data.google_project.project.project_id
      version    = "latest"
      secret     = google_secret_manager_secret.mail_recipient.secret_id
      key        = "MAIL_RECIPIENT"
    }
    service_account_email = google_service_account.function_runner.email
  }
}

resource "google_service_account" "scheduler_runner" {
  account_id = "scheduler-runner"
}

resource "google_project_iam_member" "scheduler_runner" {
  member  = "serviceAccount:${google_service_account.scheduler_runner.email}"
  project = data.google_project.project.project_id
  role    = "roles/run.invoker"
}

resource "google_cloud_scheduler_job" "scheduler" {
  name = "daily"
  http_target {
    http_method = "POST"
    uri         = google_cloudfunctions2_function.function.url
    oidc_token {
      audience              = google_cloudfunctions2_function.function.url
      service_account_email = google_service_account.scheduler_runner.email
    }
  }
  schedule  = "0 12 * * *"
  time_zone = "Asia/Tokyo"
}
