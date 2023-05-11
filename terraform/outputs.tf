output "rds_user" {
  value = aws_db_instance.db.username
}

output "rds_password" {
  value     = aws_db_instance.db.password
  sensitive = true
}

output "rds_endpoint" {
  value = aws_db_instance.db.endpoint
}

output "ecr_repository_url" {
  value = aws_ecr_repository.repository.repository_url
}

output "secret_id" {
  value = aws_secretsmanager_secret.secret.id
}

output "eks_node_role_arn" {
  value = aws_iam_role.node_role.arn
}