package concourse

import (
		"fmt"
		"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

// pipelineIDAsString converts a given numeric team ID, which is required, because Terraform resource data IDs must be
// strings.
func pipelineIDAsString(id int) string {
	return fmt.Sprintf("%d", id)
}

func resourcePipelineCreate(d *schema.ResourceData, m interface{}) error {

	name := d.Get("name").(string)
	team := d.Get("team").(string)
	paused := d.Get("paused").(bool)
	public := d.Get("public").(bool)
	config := d.Get("config").(string)

	concourse := m.(Config).Concourse().Team(team)

	// We check, if the pipeline already exists...
	pipeline, exists, err := concourse.Pipeline(name)
	if err != nil {
		return fmt.Errorf("could not fetch details of  pipeline \"%s\" prior to creation: %v", name, err)
	}
	if exists {
		return fmt.Errorf("pipeline \"%s\" does already exist in team \"%s\"", name, team)
	}

	_, _, _, err = concourse.CreateOrUpdatePipelineConfig(name, "1", []byte(config), false) // todo: see issue #3
	if err != nil {
		return fmt.Errorf("could not create pipeline config: %v", err)
	}

	// Now we check, if the pipeline has been created...
	pipeline, exists, err = concourse.Pipeline(name)
	if err != nil {
		return fmt.Errorf("could not fetch details of  pipeline \"%s\" after it's creation: %v", name, err)
	}
	if !exists {
		return fmt.Errorf("pipeline \"%s\" does not exist in team \"%s\" after an attempt to create it", name, team)
	}

	// We check if the configuration has been created.
	_, newConfig, configVersion, _, err := concourse.PipelineConfig(name)
	if err != nil {
		return fmt.Errorf("unable to read pipeline config for pipeline \"%s\" in team \"%s\" after attempting to create it: %v", name, team, err)
	}
	d.Set("config", string(newConfig))
	d.Set("config_version", configVersion)

	d.SetId(pipelineIDAsString(pipeline.ID))

	if pipeline.Paused != paused {
		var fn func(name string) (bool, error)
		if paused {
			fn = concourse.PausePipeline
		} else {
			fn = concourse.UnpausePipeline
		}
		if _, err := fn(pipeline.Name); err != nil {
			return fmt.Errorf("unable to set paused state of pipeline to %v: %v", paused, err)
		}
	}

	if pipeline.Public != public {
		var fn func(name string) (bool, error)
		if public {
			fn = concourse.ExposePipeline
		} else {
			fn = concourse.HidePipeline
		}
		if _, err := fn(pipeline.Name); err != nil {
			return fmt.Errorf("unable to set public state of pipeline to %v: %v", public, err)
		}
	}

	return nil
}

func resourcePipelineRead(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	team := d.Get("team").(string)
	name := d.Get("name").(string)

	concourse := m.(Config).Concourse().Team(team)

	pipelines, err := concourse.ListPipelines()
	if err != nil {
		return fmt.Errorf("unable to list pipelines of team \"%s\": %v", team, err)
	}

	for _, pipeline := range pipelines {
		strID := pipelineIDAsString(pipeline.ID)

		// To simplify things, we allow either the (internal) resource ID
		// or the name to be used when importing a team resource.
		if id == strID || id == pipeline.Name || (name != "" && name == pipeline.Name) {
			d.SetId(strID)
			if err := d.Set("name", pipeline.Name); err != nil {
				return err
			}
			d.Set("team", pipeline.TeamName)
			d.Set("paused", pipeline.Paused)
			d.Set("public", pipeline.Public)

			_, config, version, _, err := concourse.PipelineConfig(pipeline.Name)
			if err != nil {
				return fmt.Errorf("unable to read configuration of pipeline \"%s\": %v", pipeline.Name, err)
			}
			d.Set("config", config.String())
			d.Set("config_version", version)

			return nil
		}
	}

	// If a team with the given ID/name cannot be found, it has probably been already been deleted.
	// We will have to update the state then...
	d.SetId("")
	return nil

}

func resourcePipelineUpdate(d *schema.ResourceData, m interface{}) error {
	team := d.Get("team").(string)
	concourse := m.(Config).Concourse().Team(team)
	if d.HasChange("name") {
		var oldName , newName string
		if o, n := d.GetChange("name"); true {
			oldName = o.(string)
			newName = n.(string)
		}
		exists, err := concourse.RenamePipeline(oldName, newName)
		if err != nil {
			return fmt.Errorf("unable to rename pipeline from \"%s\" to \"%s\": %v", oldName, newName, err)
		}
		if !exists {
			return fmt.Errorf("pipeline with name \"%s\" not found", oldName)
		}
	}

	name := d.Get("name").(string)

	if d.HasChange("paused") {
		var fn func(name string) (bool, error)
		_, n := d.GetChange("paused")
		paused := n.(bool)
		if paused {
			fn = concourse.PausePipeline
		} else {
			fn = concourse.UnpausePipeline
		}
		if _, err := fn(name); err != nil {
			return fmt.Errorf("unable to set paused state of pipeline \"%s\" to %v: %v", name, paused, err)
		}
	}


	if d.HasChange("public") {
		var fn func(name string) (bool, error)
		_, n := d.GetChange("public")
		public := n.(bool)
		if public {
			fn = concourse.ExposePipeline
		} else {
			fn = concourse.HidePipeline
		}
		if _, err := fn(name); err != nil {
			return fmt.Errorf("unable to set public state of pipeline \"%s\" to %v: %v", name, public, err)
		}
	}

	if d.HasChange("config") {
		config := d.Get("config").(string)
		_, curConfig, curConfigVersionStr, _, err := concourse.PipelineConfig(name)
		if err != nil {
			return fmt.Errorf("unable to fetch configuration of pipeline \"%s\" of team \"%s\": %v", name, team, err)
		}
		curConfigVersion, err := strconv.Atoi(curConfigVersionStr)
		if err != nil {
			return fmt.Errorf("unable to convert pipeline configuration version \"%s\" to number: %v", curConfigVersionStr, err)
		}
		if string(curConfig) != config {
			_, _, _, err = concourse.CreateOrUpdatePipelineConfig(name, strconv.Itoa(curConfigVersion+1), []byte(config), false) // todo: see issue #3
			if err != nil {
				return fmt.Errorf("unable to update configuration of pipeline \"%s\" of team \"%s\" (current version: %d): %v", name, team, curConfigVersion, err)
			}
		}
	}

	return nil
}

func resourcePipelineDelete(d *schema.ResourceData, m interface{}) error {
	concourse := m.(Config).Concourse()
	team := d.Get("team").(string)
	name := d.Get("name").(string)
	_, err := concourse.Team(team).DeletePipeline(name)
	return err
}

func resourcePipelineExists(d *schema.ResourceData, m interface{}) (bool, error) {
	team := d.Get("team").(string)
	name := d.Get("name").(string)
	concourse := m.(Config).Concourse()

	// If the team does NOT exist, it makes no sense to check for pipelines of the non-existent team.
	if exists, err := teamExists(concourse, team); err != nil {
		return false, fmt.Errorf("unable to list teams: %v", err)
	} else if !exists {
		return false, nil
	}

	pipelines, err := concourse.Team(team).ListPipelines()
	if err != nil {

		return false, fmt.Errorf("unable to list pipelines: %v", err)
	}
	for _, pipeline := range pipelines {
		if pipeline.Name == name && pipeline.TeamName == team {
			return true, nil
		}
	}

	return false, nil

}

func resourcePipelineState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	nameOrID := d.Id()
	if err := resourcePipelineRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("pipeline with ID or name %s not found", nameOrID)
	}
	return []*schema.ResourceData{d}, nil

}

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Exists: resourcePipelineExists,
		Schema: map[string]*schema.Schema{
			"team": {
				Description: "Team name",
				Type: schema.TypeString,
				Required: true,
				ForceNew: true, // Movement of pipelines between teams is not supported at the moment
			},
			"name": {
				Description: "Pipeline name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"paused": {
				Description: "Paused",
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"public": {
				Description: "Public",
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"config": {
				Description: "Pipeline configuration YAML",
				Type: schema.TypeString,
				Required: true,
			},
			"config_version": {
				Description: "Pipeline configuration version",
				Type: schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourcePipelineState,
		},
	}
}
