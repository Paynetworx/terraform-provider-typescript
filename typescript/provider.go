package typescript

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"typescript_file": resourceTypescriptFile(),
			"typescript_node_modules": resourceTypescriptNodeModules(),
		},
	}
}
