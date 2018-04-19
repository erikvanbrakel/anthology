resource "aws_ecs_service" "registry-service" {
  name = "terraform-registry"

  task_definition = "${aws_ecs_task_definition.registry-task.id}"
  cluster         = "${var.ecs_cluster}"
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = ["${aws_security_group.registry.id}"]
    subnets         = ["${var.private_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_lb_target_group.registry.arn}"
    container_name   = "registry"
    container_port   = 8082
  }

  depends_on = ["aws_lb_listener.https"]
}

resource "aws_ecs_task_definition" "registry-task" {
  family                   = "terraform-registry"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512

  task_role_arn      = "${aws_iam_role.task-role.arn}"
  execution_role_arn = "${aws_iam_role.execution-role.arn}"

  container_definitions = "${data.template_file.container-definitions.rendered}"
}

data "aws_region" "current" {}

data "template_file" "container-definitions" {
  template = "${file("${path.module}/templates/container-definitions.json")}"

  vars {
    BUCKET_NAME    = "${var.bucket}"
    REGISTRY_IMAGE = "erikvanbrakel/terraform-registry:latest"
    LOG_GROUP      = "${aws_cloudwatch_log_group.log-group.name}"
    LOG_REGION     = "${data.aws_region.current.name}"
  }
}

resource "aws_cloudwatch_log_group" "log-group" {
  name              = "terraform-registry"
  retention_in_days = 7
}
