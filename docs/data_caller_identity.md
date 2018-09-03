## Data Source: concourse_caller_identity

Use this data source to get access to information about the current connection
to the Concourse ATC/web server.

### Example Usage

```hcl
data "concourse_caller_identity" "current" {}

output "user_name" {
  value = "${concourse_caller_identity.current.user_name}"
}

output "is_admin" {
  value = "${concourse_caller_identity.current.is_admin}"
}
```

### Argument Reference

There are no arguments available for this data source.

### Attribute Reference

* `user_name` - User name used for the current connection to the Concourse ATC/web server.
* `is_admin` - Boolean flag that indicates if the current user has admin privileges.
