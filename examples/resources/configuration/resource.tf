
resource "ravel_configuration" "example" {
  name = "test"

  scope = {
    terraform = "testing"
  }

  labels = {
    minDominoVersion = "005.008.000"
  }

  definition = jsonencode({
    "email_notifications" : {
      "enabled" : true,
      "server" : "smtp.customer.org",
      "port" : 465,
      "enable_ssl" : true,
      "from_address" : "domino@customer.org",
      "authentication" : {
        "username" : "test",
        "password" : "123"
      }
    }
  })
}

output "id" {
  value = ravel_configuration.example.id
}

output "version" {
  value = ravel_configuration.example.version
}
