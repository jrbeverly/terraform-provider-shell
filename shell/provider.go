package shell

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"working_directory": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PWD", nil),
				Description: "The working directory where to run.",
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"shell": resourceExternal(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cmd := CmdRunner{
		TemporaryDirectory: os.TempDir(),
		RetryMaximum:       5,
	}

	config := Config{
		WorkingDirectory: d.Get("working_directory").(string),
		Runner:           cmd,
		Variables:        d.Get("variables").(map[string]interface{}),
	}

	return &config, nil
}
