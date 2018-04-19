variable "tld" {
  type        = "string"
  description = "TLD under which the registry will be hosted."
}

variable "storage_bucket" {
  type        = "string"
  description = "This bucket is used to store the terraform modules."
}

data "aws_route53_zone" "zone" {
  name = "${var.tld}"
}

resource "aws_s3_bucket" "modules" {
  bucket = "${var.storage_bucket}"

  force_destroy = true
}

resource "aws_route53_record" "cname" {
  ttl     = 60
  name    = "registry"
  type    = "CNAME"
  zone_id = "${data.aws_route53_zone.zone.id}"
  records = ["${module.registry.dns_name}"]
}

resource "aws_acm_certificate" "certificate" {
  domain_name       = "registry.${var.tld}"
  validation_method = "DNS"
}

resource "aws_route53_record" "validation" {
  name    = "${aws_acm_certificate.certificate.domain_validation_options.0.resource_record_name}"
  type    = "${aws_acm_certificate.certificate.domain_validation_options.0.resource_record_type}"
  zone_id = "${data.aws_route53_zone.zone.id}"
  records = ["${aws_acm_certificate.certificate.domain_validation_options.0.resource_record_value}"]
  ttl     = 60
}

resource "aws_ecs_cluster" "cluster" {
  name = "registry-cluster"
}

module "networking" {
  source     = "./networking"
  cidr_block = "192.168.0.0/20"
}

module "registry" {
  source = "./registry"

  public_subnet_ids  = "${module.networking.public_subnet_ids}"
  private_subnet_ids = "${module.networking.private_subnets_ids}"
  vpc_id             = "${module.networking.vpc_id}"

  certificate_arn = "${aws_acm_certificate.certificate.arn}"
  bucket          = "${aws_s3_bucket.modules.bucket}"
  ecs_cluster     = "${aws_ecs_cluster.cluster.name}"
}
