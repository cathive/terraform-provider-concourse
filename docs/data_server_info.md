## Data Source: concourse_server_info

Use this data source to get access to information about the Concourse ATC/web
server currently connected to.

### Example Usage

```hcl
data "concourse_server_info" "current" {}

output "version" {
  value = "${concourse_server_info.current.version}"
}

output "worker_version" {
  value = "${concourse_server_info.current.worker_version}"
}
```

### Argument Reference

There are no arguments available for this data source.

### Attribute Reference

* `version` - The version of the Concourse ATC/web server currently connected to.
* `worker_version` - The worker version that the Concourse ATC/web server is compatible with.
