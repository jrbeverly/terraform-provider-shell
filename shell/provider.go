package shell

import (
	"github.com/Brightspace/terraform-provider-shell/shell/api"
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
			"max_retries": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "This is the maximum number of times an API call is retried, in the case where requests are being throttled or experiencing transient failures.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"shell": resourceExternal(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	retry_maximum := 5
	if v, ok := d.GetOk("max_retries"); ok {
		retry_maximum = v.(int)
	}

	cmd := api.CmdRunner{
		TemporaryDirectory: os.TempDir(),
		RetryMaximum:       retry_maximum,
	}

	config := Config{
		WorkingDirectory: d.Get("working_directory").(string),
		Runner:           cmd,
		Variables:        d.Get("variables").(map[string]interface{}),
	}

	return &config, nil
}
