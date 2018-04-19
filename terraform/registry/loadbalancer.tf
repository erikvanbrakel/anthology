resource "aws_lb" "loadbalancer" {
  name = "terraform-registry"

  subnets         = ["${var.public_subnet_ids}"]
  security_groups = ["${aws_security_group.loadbalancer.id}"]
}

resource "aws_lb_target_group" "registry" {
  name = "registry"

  vpc_id      = "${var.vpc_id}"
  target_type = "ip"
  protocol    = "HTTP"
  port        = 80

  health_check {
    port                = "traffic-port"
    path                = "/.well-known/terraform.json"
    healthy_threshold   = 5
    unhealthy_threshold = 2
    interval            = 30
    timeout             = 5
  }
}

resource "aws_lb_listener" "https" {

  load_balancer_arn = "${aws_lb.loadbalancer.arn}"
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = "${var.certificate_arn}"

  default_action {
    target_group_arn = "${aws_lb_target_group.registry.id}"
    type             = "forward"
  }
}
