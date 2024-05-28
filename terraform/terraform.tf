provider "google" {
  project = "ur-available-room"
  region  = "asia-northeast1"
}

terraform {
  backend "gcs" {
    bucket = "ur-available-room-tf-backend"
    prefix = "terraform.tfstate"
  }
}
