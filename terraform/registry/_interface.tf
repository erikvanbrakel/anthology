variable "vpc_id" {
  description = "Target VPC ID. Required for creating security groups."
  type = "string"
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs. The Application Loadbalancer will be deployed into these."
  type = "list"
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs. The registry will be deployed into one of these."
  type = "list"
}

variable "certificate_arn" {
  description = "ARN of the certificate to use for TLS on the load balancer."
  type = "string"
}

variable "bucket" {
  description = "Bucket used to store the modules."
  type = "string"
}

variable "ecs_cluster" {
  description = "Cluster to run the registry service in."
  type = "string"
}

output "dns_name" {
  value = "${aws_lb.loadbalancer.dns_name}"
}
