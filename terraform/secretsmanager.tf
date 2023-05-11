resource "aws_secretsmanager_secret" "secret" {
  name                    = var.app
  recovery_window_in_days = 0
}