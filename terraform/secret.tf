resource "google_secret_manager_secret" "mailgun_api_key" {
  replication {
    auto {}
  }
  secret_id = "mailgun-api-key"
}

resource "google_secret_manager_secret_version" "secret" {
  secret      = google_secret_manager_secret.mailgun_api_key.id
  secret_data = "mailgun-api-key"
}

resource "google_secret_manager_secret" "mailgun_domain" {
  replication {
    auto {}
  }
  secret_id = "mailgun-domain"
}

resource "google_secret_manager_secret_version" "mailgun_domain" {
  secret      = google_secret_manager_secret.mailgun_domain.id
  secret_data = "mailgun-domain"
}

resource "google_secret_manager_secret" "mail_recipient" {
  replication {
    auto {}
  }
  secret_id = "mailgun-recipient"
}

resource "google_secret_manager_secret_version" "mail_recipient" {
  secret      = google_secret_manager_secret.mail_recipient.id
  secret_data = "mail-recipient"
}
