resource "local_file" "foo" {
    content     = "foo!"
    filename = "${path.module}/foo.bar"
}

module "s3" {
    source = "../../modules/s3"
    name = "my-bucket"
}

module "user" {
    source = "../../modules/iam-user"
    name = "george"
}
