---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ravel Provider"
subcategory: ""
description: |-
  
---

# ravel Provider



## Example Usage

```terraform
provider "ravel" {
  url = "https://domino.ai/ravel"

  token = "ABC-123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `token` (String, Sensitive) The access token for API operations. Can be defined from env var RAVEL_TOKEN
- `url` (String) Host URL for Ravel. Can be defined from env var RAVEL_URL