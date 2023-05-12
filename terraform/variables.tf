variable "region" {
  type    = string
  default = "ue1"
}

variable "environment" {
  type    = string
  default = "p"
}

variable "app" {
  type    = string
  default = "golang-bank"
}

variable "tags" {
  type    = map(string)
  default = {
    Application = "Golang Bank"
    Environment = "Production"
  }
}