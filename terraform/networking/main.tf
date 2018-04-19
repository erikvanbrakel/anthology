data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_availability_zones" "zones" {}

resource "aws_vpc" "vpc" {
  cidr_block = "${var.cidr_block}"

  tags {
    Name = "registry-vpc-${data.aws_region.current.name}"
  }
}

resource "aws_internet_gateway" "internet-gateway" {
  vpc_id = "${aws_vpc.vpc.id}"

  tags {
    Name = "${aws_vpc.vpc.tags["Name"]}-igw"
  }
}

resource "aws_subnet" "public" {
  count = 3

  vpc_id            = "${aws_vpc.vpc.id}"
  cidr_block        = "${cidrsubnet(aws_vpc.vpc.cidr_block, 8, count.index)}"
  availability_zone = "${element(data.aws_availability_zones.zones.names, count.index)}"

  tags {
    Name = "${aws_vpc.vpc.tags["Name"]}-public${count.index}"
  }
}

resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.vpc.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.internet-gateway.id}"
  }
}

resource "aws_route_table_association" "public" {
  count = 3

  route_table_id = "${aws_route_table.public.id}"
  subnet_id      = "${element(aws_subnet.public.*.id, count.index)}"
}

resource "aws_subnet" "private" {
  count = 3

  vpc_id            = "${aws_vpc.vpc.id}"
  cidr_block        = "${cidrsubnet(aws_vpc.vpc.cidr_block, 8, count.index + 3)}"
  availability_zone = "${element(data.aws_availability_zones.zones.names, count.index)}"

  tags {
    Name = "${aws_vpc.vpc.tags["Name"]}-private${count.index}"
  }
}

resource "aws_eip" "eip" {
  tags {
    Name = "${aws_vpc.vpc.tags["Name"]}-nat"
  }
  depends_on = ["aws_internet_gateway.internet-gateway"]
}

resource "aws_nat_gateway" "nat" {
  allocation_id = "${aws_eip.eip.id}"
  subnet_id     = "${element(aws_subnet.public.*.id,0)}"

  tags {
    Name = "${aws_vpc.vpc.tags["Name"]}-nat"
  }

}

resource "aws_route_table" "private" {
  vpc_id = "${aws_vpc.vpc.id}"

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = "${aws_nat_gateway.nat.id}"
  }
}

resource "aws_route_table_association" "private" {
  count = 3

  route_table_id = "${aws_route_table.private.id}"
  subnet_id      = "${element(aws_subnet.private.*.id,count.index)}"
}
