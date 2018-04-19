variable "cidr_block" {
  type        = "string"
  description = "CIDR block of the VPC. Subnets will be +8. e.g. /16 results in /16+8 => /24 subnets."
}

output "public_subnet_ids" {
  value = "${aws_subnet.public.*.id}"
}

output "private_subnets_ids" {
  value = "${aws_subnet.private.*.id}"
}

output "vpc_id" {
  value = "${aws_vpc.vpc.id}"
}
