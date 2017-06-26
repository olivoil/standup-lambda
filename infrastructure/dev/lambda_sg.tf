resource "aws_security_group" "lambda_security_group" {
  name = "${var.project}-${var.apex_environment}-lambda-security-group"
  description = "Security Group for lambda functions"
  vpc_id = "${module.vpc.vpc_id}"

  // allows traffic from the SG itself for tcp
  ingress {
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    self = true
  }

  // allows traffic from the SG itself for udp
  ingress {
    from_port = 0
    to_port = 65535
    protocol = "udp"
    self = true
  }

  // allows all traffic out
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}