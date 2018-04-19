resource "aws_security_group" "loadbalancer" {
  vpc_id = "${var.vpc_id}"
  name   = "public-loadbalancer"

  ingress {
    from_port   = 443
    protocol    = "TCP"
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS from everywhere"
  }

  egress {
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    description = "Outgoing connections to everywhere"
  }
}

resource "aws_security_group" "registry" {
  vpc_id = "${var.vpc_id}"
  name   = "terraform-registry"

  ingress {
    from_port       = 0
    protocol        = "TCP"
    to_port         = 65535
    security_groups = ["${aws_security_group.loadbalancer.id}"]
    description     = "All TCP traffic from the loadbalancer group"
  }

  egress {
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    description = "Outgoing connections to everywhere"
  }
}
