# ---------------------------------------------------------------------------------------------------------------------
# BASIC TERRAFORM EXAMPLE
# See test/terraform_aws_example.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "example" {
  template = "${var.example}"
}