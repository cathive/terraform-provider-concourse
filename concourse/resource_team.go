package concourse

import (
	"fmt"
	"github.com/concourse/concourse/atc"

	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/hashicorp/terraform/helper/schema"
)

// teamIDAsString converts a given numeric team ID, which is required, because Terraform resource data IDs must be
// strings.
func teamIDAsString(id int) string {
	return fmt.Sprintf("%d", id)
}

func teamExists(concourse concourse.Client, name string) (bool, error) {
	teams, err := concourse.ListTeams()
	if err != nil {
		return false, fmt.Errorf("unable to list teams: %v", err)
	}
	for _, team := range teams {
		if team.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func resourceTeamCreate(d *schema.ResourceData, m interface{}) error {
	concourse := m.(Config).Concourse()

	name := d.Get("name").(string)

	team := atc.Team{
		Name: name,
		Auth: atc.TeamAuth{},
	}
	for _, auth := range(d.Get("auth").([]interface{})) {
		authElem := auth.(struct{
			scope string
			users []string
			groups []string
		})
		team.Auth[authElem.scope] = map[string][]string{
			"users": authElem.users,
			"groups": authElem.groups,
		}
	}

	team, created, updated, err := concourse.Team(name).CreateOrUpdate(team)
	if err != nil {
		return fmt.Errorf("could not create team: %v", err)
	}
	if !created && !updated {
		return fmt.Errorf("could not create/update team %s: neither 'created' nor 'updated' was set to true", name)
	}
	d.SetId(teamIDAsString(team.ID))
	return nil
}

func resourceTeamRead(d *schema.ResourceData, m interface{}) error {
	concourse := m.(Config).Concourse()

	id := d.Id()
	name := d.Get("name").(string)

	teams, err := concourse.ListTeams()
	if err != nil {
		return fmt.Errorf("unable to list teams: %v", err)
	}

	for _, team := range teams {
		strID := teamIDAsString(team.ID)

		// To simplify things, we allow either the (internal) resource ID
		// or the name to be used when importing a team resource.
		if id == strID || id == team.Name || (name != "" && name == team.Name) {
			d.SetId(strID)
			if err := d.Set("name", team.Name); err != nil {
				return err
			}
			return nil
		}
	}

	// If a team with the given ID/name cannot be found, it has probably been already been deleted.
	// We will have to update the state then...
	d.SetId("")
	return nil

}

func resourceTeamUpdate(d *schema.ResourceData, m interface{}) error {
	concourse := m.(Config).Concourse()
	newName := ""
	if d.HasChange("name") {
		newName = d.Get("name").(string)
	}
	teams, err := concourse.ListTeams()
	id := d.Id()
	if err != nil {
		return fmt.Errorf("unable to list teams: %v", err)
	}
	for _, team := range teams {
		if id == teamIDAsString(team.ID) {
			oldName := team.Name
			if newName != "" && oldName != newName {
				_, err := concourse.Team(team.Name).RenameTeam(team.Name, newName)
				if err != nil {
					return fmt.Errorf("unable to rename team %s to %s: %v", oldName, newName, err)
				}
				// We store the new name because we might change other properties of the team as well.
				team.Name = newName
			}
			if d.HasChange("auth") {
				team.Auth = atc.TeamAuth{}

				for _, auth := range(d.Get("auth").([]interface{})) {
					authElem := auth.(struct{
						scope string
						users []string
						groups []string
					})
					team.Auth[authElem.scope] = map[string][]string{
						"users": authElem.users,
						"groups": authElem.groups,
					}
				}

				team, _, _, err = concourse.Team(newName).CreateOrUpdate(team)
				if err != nil {
					return fmt.Errorf("unable to update team %s: %v", newName, err)
				}
			}
			return nil
		}
	}
	return fmt.Errorf("team with ID %s not found", d.Id())

}

func resourceTeamDelete(d *schema.ResourceData, m interface{}) error {
	concourse := m.(Config).Concourse()
	name := d.Get("name").(string)
	return concourse.Team(name).DestroyTeam(name)
}

func resourceTeamExists(d *schema.ResourceData, m interface{}) (bool, error) {
	id := d.Id()
	concourse := m.(Config).Concourse()

	teams, err := concourse.ListTeams()
	if err != nil {
		return false, fmt.Errorf("unable to list teams: %v", err)
	}
	for _, team := range teams {
		teamID := teamIDAsString(team.ID)
		if teamID == id {
			return true, nil
		}
	}
	return false, nil
}

func resourceTeamState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	nameOrID := d.Id()
	if err := resourceTeamRead(d, m); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("team with ID or name %s not found", nameOrID)
	}
	return []*schema.ResourceData{d}, nil

}

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,
		Exists: resourceTeamExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Team name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auth": {
				Description: "Access and authorization settings for users and groups",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Description: "Role within the team",
							Type: schema.TypeString,
							Required: true,
						},
						"users": {
							Type: schema.TypeList,
							Elem: schema.TypeString,
							Optional: true,
							PromoteSingle: true,
						},
						"groups": {
							Type: schema.TypeList,
							Elem: schema.TypeString,
							Optional: true,
							PromoteSingle: true,
						},
					},
				},
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceTeamState,
		},
	}
}
