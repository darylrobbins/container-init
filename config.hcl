env "foo" {
  value = "bar"
}

service {
  process {
    command = ["/usr/local/bin/awesome-app", "server"]
  }
}