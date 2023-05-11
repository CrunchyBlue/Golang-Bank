data "aws_iam_policy_document" "cluster_policy_document" {
  version = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = [
      "sts:AssumeRole"
    ]
    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "cluster_role" {
  name                  = "AWSEKSClusterRole"
  assume_role_policy    = data.aws_iam_policy_document.cluster_policy_document.json
  description           = "Allows access to other AWS service resources that are required to operate clusters managed by EKS."
  force_detach_policies = false
  managed_policy_arns   = ["arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"]
  max_session_duration  = 3600
  path                  = "/"
}

resource "aws_iam_role_policy_attachment" "AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster_role.name
}