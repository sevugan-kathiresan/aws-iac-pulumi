package tags

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type Tags struct {
	Creator string
	Team    string
	Service string
	Env     string
}

// Function
// Purpose: Instantiates the Tags struct
func NewTags(creator, team, service, env string) Tags {
	return Tags{
		Creator: creator,
		Team:    team,
		Service: service,
		Env:     env,
	}
}

// Method: Associated with the struct Tags
// Purpsose: to convert the Tags struct to a map so that we can use it with Pulumi AWS resources
func (t Tags) ToPulumiMap() pulumi.StringMap {
	return pulumi.StringMap{
		"Creator":    pulumi.String(t.Creator),
		"Team":       pulumi.String(t.Team),
		"Service":    pulumi.String(t.Service),
		"Env":        pulumi.String(t.Env),
		"Managed_by": pulumi.String("Pulumi"),
	}
}
