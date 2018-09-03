## Data Source: concourse_team

Use this data source to get access to information about a team.

### Example Usage

```hcl
data "concourse_team" "main" {
  name = "main"
}

output "numeric_id" {
  value = "${concourse_team.main.id}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - Name of the team (required).

### Attribute Reference

in addition to all arguments above, the following attributes are exported:

* `id` - Numeric unique ID of the team.
