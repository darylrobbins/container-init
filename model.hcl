variable "log_level" {
  type = string
  allowed_values = ["DEBUG", "INFO", "WARN", "ERROR"]
  default = "INFO"
}

secret "db_secret" {
  value = lookup(env.DB_SECRET_ARN)
  format = "json"
}

env "DD_AGENT_HOST" {
  value = file("/host/local_ipv4")
}

file "log4j" {
  path = "/app/log4j.xml"
  content = templatefile("${path.module}/log4j.xml", {
    log_level = var.log_level
  })
  on_change = "signal(SIG_HUP)"
}

file "baseline" {
  path = "/app/baseline.yaml"
  content = file("s3://${env.CONFIG_BUCKET}/baseline.yaml")
}

file "app_config" {
  path = "/app/app_config.json"
  content = jsonencode({
    db_user = secret.db_secret.user,
    db_password = secret.db_secret.password
    stuff_dir = directory.stuff.path
  })
  on_change = "restart"
}

directory "stuff" {
  path = "/app/stuff"
  source = "s3://${env.CONFIG_BUCKET}/stuff"
  exclude = "**/*.tmp"
}

service {
  process {
    command = ["/usr/local/bin/awesome-app", "server", file.app_config.path]
  }
}