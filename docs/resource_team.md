## concourse_team

### Example Usage

```hcl
resource "concourse_team" "team_a" {
  name = "team-a"
}
```

### Argument Reference

The following arguments are supported:

* `name` - Name of the team.

### Attributes Reference

in addition to all arguments above, the following attributes are exported:

* `id` - Numeric unique ID of the team.

### Import

Teams can be imported using either their `name` or their `numeric ID`, e.g.:

```sh
$ terraform import concourse_team.my_team my-team
```