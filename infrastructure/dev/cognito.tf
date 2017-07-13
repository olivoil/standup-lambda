
resource "aws_cognito_identity_pool" "dev" {
  identity_pool_name               = "${var.project}_${var.apex_environment}_identity_pool"
  allow_unauthenticated_identities = false

  supported_login_providers {
    "accounts.google.com" = "218118279936-leo41neq68o0ltcq471e44d09tn84pkl.apps.googleusercontent.com"
  }
}

// unauthenticated role
resource "aws_iam_role" "unauthenticated_cognito_role" {
  name = "${var.project}-${var.environment}-unauthenticated-cognito-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.dev.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "unauthenticated"
        }
      }
    }
  ]
}
EOF
}

// authenticated role
resource "aws_iam_role" "authenticated_cognito_role" {
  name = "${var.project}-${var.environment}-authenticated-cognito-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.dev.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "accounts.google.com"
        }
      }
    }
  ]
}
EOF
}
