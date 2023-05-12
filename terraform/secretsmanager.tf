resource "random_password" "password" {
  length  = 32
  special = false
}

resource "random_password" "symmetric_key" {
  length  = 32
  special = false
}

resource "aws_secretsmanager_secret" "secret" {
  name                    = var.app
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "secret_val" {
  secret_id     = aws_secretsmanager_secret.secret.id
  secret_string = jsonencode({
    "DB_DRIVER" : "postgres",
    "DB_SOURCE" : "postgresql://${aws_db_instance.db.username}:${random_password.password.result}@${aws_db_instance.db.endpoint}",
    "SERVER_ADDRESS" : "0.0.0.0:8080",
    "ACCESS_TOKEN_SYMMETRIC_KEY" : random_password.symmetric_key.result,
    "ACCESS_TOKEN_DURATION" : "15m"
  })
}