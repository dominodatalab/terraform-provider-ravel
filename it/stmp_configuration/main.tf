resource "ravel_configuration" "stmp_configuration" {
  name = "domino-cloud-smtp-configuration-test"

  labels = {
    minSchemaVersion = "001.000"
  }
  scope = {
    type                 = "configuration"
    category             = "fleetcommand-configuration-manager"
    fleetcommand_account = "ldebello-account"
  }

  schema = {
    name    = "email"
    version = "1.0.0"

    scope = {
      type     = "schema"
      source   = "domino/release"
      category = "fleetcommand-configuration-manager"
    }
  }

  definition = jsonencode({
    "email_notifications" : {
      "enabled" : true,
      "server" : "email-smtp.us-east-1.amazonaws.com",
      "port" : 465,
      "enable_ssl" : true,
      "from_address" : "cloud-support@dominodatalab.com",
      "authentication" : {
        "username" : "user",
        "password" : "password"
      }
    }
  })
}
