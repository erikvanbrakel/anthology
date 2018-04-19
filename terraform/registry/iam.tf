resource "aws_iam_role" "task-role" {
  name = "registry-task-role"

  assume_role_policy = "${data.aws_iam_policy_document.ecs-assume-role.json}"
}

resource "aws_iam_role" "execution-role" {
  name = "registry-execution-role"

  assume_role_policy = "${data.aws_iam_policy_document.ecs-assume-role.json}"
}

data "aws_iam_policy_document" "ecs-assume-role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "execution" {
  role       = "${aws_iam_role.execution-role.id}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "task-s3" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
  role       = "${aws_iam_role.task-role.id}"
}
