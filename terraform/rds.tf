resource "aws_db_instance" "db" {
  allocated_storage            = 20
  copy_tags_to_snapshot        = true
  identifier                   = var.app
  db_name                      = "bank"
  engine                       = "postgres"
  engine_version               = "15.2"
  instance_class               = "db.t3.micro"
  license_model                = "postgresql-license"
  manage_master_user_password  = true
  multi_az                     = false
  network_type                 = "IPV4"
  option_group_name            = "default:postgres-15"
  parameter_group_name         = "default.postgres15"
  performance_insights_enabled = true
  publicly_accessible          = true
  skip_final_snapshot          = true
  storage_encrypted            = true
  username                     = "root"
}