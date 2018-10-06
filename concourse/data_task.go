package concourse

import "github.com/hashicorp/terraform/helper/schema"

func dataTask() *schema.Resource {
	return &schema.Resource{
		//Read:   resourceTeamRead,
		//Exists: resourceTeamExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Team name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config": {
				Description: "Concourse Task configuration block",
				Type:        schema.TypeList, // Todo: remove this? only allow one config block
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// https://concourse-ci.org/tasks.html#task-platform
						"platform": {
							Description: "Platorm the container will run on",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "linux",
						},

						// https://concourse-ci.org/tasks.html#task-image-resource
						"image_resource": {
							//Type: schema.TypeList,
							Optional: true, // Not sure why optional - what happens if run without one?
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Image Type",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "docker-image",
									},
									"source": {
										Description: "Image Source",
										Type:        schema.TypeList,
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"repository": {
													Type:     schema.TypeString,
													Required: true,
												},
												"tag": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "latest",
												},
												// Todo: "params": {}
												// Todo: "version": {}
											},
										},
									},
								},
							},
						},

						// https://concourse-ci.org/tasks.html#task-params
						"params": {
							Description: "Parameters to pass in as environment variables",
							Type:        schema.TypeMap,
							Optional:    true,
						},

						// https://concourse-ci.org/tasks.html#task-inputs
						"input": {
							Description: "Input from another task or resource",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Name of the input",
										Type:        schema.TypeString,
										Required:    true,
									},
									"path": {
										Description: "Directory path of input (defaults to name)",
										Type:        schema.TypeString,
										Optional:    true,
									},
									// Todo: "optional": {}
								},
							},
						},

						// https://concourse-ci.org/tasks.html#task-outputs
						"output": {
							Description: "Output to another task or resource",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Name of the output",
										Type:        schema.TypeString,
										Required:    true,
									},
									"path": {
										Description: "Directory path of output (defaults to name)",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},

						// https://concourse-ci.org/tasks.html#task-run
						"run":     {
							Description: "Execute commands and scripts",
							Type:        schema.TypeList, // Todo: remove this? Only allow one run block
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Name of the output",
										Type:        schema.TypeString,
										Required:    true,
									},
									"path": {
										Description: "Directory path of output (defaults to name)",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},

						// https://concourse-ci.org/tasks.html#task-caches
						// Todo: "caches": {}
					},
				},
			},
			// Todo: "on_success": {},
			// Todo: "on_failure": {},
			// Todo: "ensure": {},
		},
	}
}
