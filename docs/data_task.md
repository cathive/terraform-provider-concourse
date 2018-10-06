## Data Source: concourse_task

Use this data source to generate a task definition to bind to a concourse_job data source.

### Example Usage

```hcl
variable "script" {
  default = <<EOF
env
find . -type f
echo "Hi"
EOF
}

data "concourse_task" "task1" {
  name = "task1"

  config {
    platform = "linux"

    image_resource {
      type = "docker-image"

      source {
        repository = "hashicorp/terraform"
        tag        = "light"
      }
    }

    run {
      path = "sh"
      args = ["-c", "${var.script}"]
      user = "python"
      dir  = "one"
    }

    params = "${map(
        "A", "1",
        "B", "2",
        "C", "3"
    )}"

    input {
      name = "input_1"
      path = "one"
    }

    input {
      name = "input_2"
    }

    output {
      name = "output_1"
      path = "one"
    }

    output {
      name = "output_1"
      path = "one"
    }
  }

  on_success {
    // Todo
  }

  on_failure {
    // Todo
  }

  ensure {
    // Todo
  }
}

```

### Argument Reference

The following arguments are supported:

* `name` - Name of the team (required).

### Attribute Reference

in addition to all arguments above, the following attributes are exported:

* `id` - Numeric unique ID of the team.
