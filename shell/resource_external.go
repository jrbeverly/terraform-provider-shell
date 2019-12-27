package shell

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceExternal() *schema.Resource {
	return &schema.Resource{
		Create: resourceExternalCreate,
		Update: resourceExternalUpdate,
		Read:   resourceExternalRead,
		Delete: resourceExternalDelete,

		Schema: map[string]*schema.Schema{
			"create": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"update": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delete": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"read": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"query": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceExternalCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	env_vars := make(map[string]interface{})
	programI := d.Get("create").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	for k, v := range query {
		env_vars[k] = v
	}
	for k, v := range config.Variables {
		env_vars[k] = v
	}

	result, err := runShellCommand(programI, workingDir, env_vars, d.Id())
	if err != nil {
		return fmt.Errorf("create: %s", err)
	}

	d.Set("result", result)
	d.SetId(result["id"].(string))
	log.Printf("[INFO] Created generic resource: %s", d.Id())

	return nil
}

func resourceExternalRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	env_vars := make(map[string]interface{})
	programI := d.Get("read").([]interface{})
	log.Printf("[INFO] Number of command args [%d]", len(programI))
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	for k, v := range query {
		env_vars[k] = v
	}
	for k, v := range config.Variables {
		env_vars[k] = v
	}

	result, err := runShellCommand(programI, workingDir, env_vars, d.Id())
	if err != nil {
		log.Printf("[INFO] Error occurred while retrieving resource %s", d.Id())
		d.SetId("")
		return fmt.Errorf("read: %s", err)
	}

	if result["id"] == "" {
		log.Printf("[INFO] Resource could not be found %s", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("result", result)
	log.Printf("[INFO] Created generic resource: %s", d.Id())

	return nil
}

func resourceExternalUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	env_vars := make(map[string]interface{})
	programI := d.Get("update").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	for k, v := range query {
		env_vars[k] = v
	}
	for k, v := range config.Variables {
		env_vars[k] = v
	}

	result, err := runShellCommand(programI, workingDir, env_vars, d.Id())
	if err != nil {
		return fmt.Errorf("update: %s", err)
	}

	d.Set("result", result)
	log.Printf("[INFO] Updated resource: %s", d.Id())

	return nil
}

func resourceExternalDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	env_vars := make(map[string]interface{})
	programI := d.Get("delete").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	for k, v := range query {
		env_vars[k] = v
	}
	for k, v := range config.Variables {
		env_vars[k] = v
	}

	result, err := runShellCommand(programI, workingDir, env_vars, d.Id())
	if err != nil {
		return fmt.Errorf("delete: %s", err)
	}

	log.Printf("[INFO] Deleted resource: %s", d.Id())
	d.Set("result", result)
	d.SetId("")

	return nil
}
