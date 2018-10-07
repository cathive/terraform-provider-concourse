## Data Source: concourse_task

Use this data source to generate a task definition to bind to a concourse_job data source.

### Example Usage

```hcl

data "concourse_task" "task1" {
  name = "task1"

  config {
    image_resource {
      source {
        repository = "alpine"
        tag        = "3.8"
      }
    }

    run {
      path = "./one/myscript.sh"
    }

    input {
      name = "input_1"
      path = "one"
    }

    output {
      name = "slack-attachments"
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

* `name` (Required) - Name of the team.
* `config` (Required) - Config block containing the task configuration.
  * `platform` (Optional) - Platform to be run on (defaults to `linux`).
  * `image_resource` (Optional) - Config block containing the image_resource configuration.
    * `type` (Optional) - Image_resource type (defaults to `docker-image`)
    * `source` (Required) - Image Source configuration block.
      * `repository` (Required) - Standard definition for a container image and optionally the registry, eg: `<registry>/<repository>`.
      * `tag` (Optional) - Image tag (defaults to `latest`)
  * `params` (Optional) - A map of variable name to value to be passed through to the container as environment variables.
  * `input` (Optional) - Input from another task or resource. This can defined more than once.
    * `name` (Required) - Name of the input. This appears as a directory in the container.
    * `path` (Optional) - Directory path of input if it should be different to the name.
  * `output` (Optional) - Input from another task or resource. This can defined more than once.
    * `name` (Required) - Name of the output. This appears as a directory in the container.
    * `path` (Optional) - Directory path of output if it should be different to the name.
  * `run` (Required) - Configuration block for commands and scripts to be executed.
    * `path` (Required) - The command to execute.
    * `args` (Optional) - List of arguments to be passed to the command.
    * `user` (Optional) - Runtime user (defaults to `root` or whatever is specified by Garden)
    * `dir` (Optional) - Relative working directory to change to before executing the command.
* `on_success` (Optional) - Steps to run on successful task completion.
  * _Todo_
* `on_failure` (Optional) - Steps to run on task failure.
  * _Todo_
* `ensure` (Optional) - Steps to run at the end of a task regardless of success or failure.
  * _Todo_

### Attribute Reference

in addition to all arguments above, the following attributes are exported:

* `yaml` - The above arguments serialized as a standard YAML task definition.

### Full Example:
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