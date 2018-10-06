output "instance_id" {
  value = "${google_compute_instance.example.id}"
}

output "bucket_url" {
  value = "${google_storage_bucket.example_bucket.url}"
}
