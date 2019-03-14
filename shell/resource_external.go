package shell

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

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

func convertToEnvVars(args map[string]interface{}) []string {
	i := 0
	vars := make([]string, len(args))
	for key, val := range args {
		vars[i] = fmt.Sprintf("%s=%s", key, val.(string))
		i++
	}
	vars = append(vars, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	return vars
}

func runCommand(programI []interface{}, workingDir string, query map[string]interface{}, id string) (map[string]interface{}, error) {
	log.Printf("[INFO] Number of command args [%d]", len(programI))
	program := make([]string, len(programI))
	for i, vI := range programI {
		log.Printf("[INFO] Program [%d]: %s", i, vI.(string))
		program[i] = vI.(string)
	}
	if len(program) == 0 {
		return nil, fmt.Errorf("No command has been provided")
	}

	env := convertToEnvVars(query)

	cmd := exec.Command(program[0], program[1:]...)
	cmd.Dir = workingDir
	cmd.Env = append(env, fmt.Sprintf("TF_ID=%s", id))

	resultJson, err := cmd.Output()
	log.Printf("[INFO] result %s", resultJson)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.Stderr != nil && len(exitErr.Stderr) > 0 {
				return nil, fmt.Errorf("failed to execute %q: %s", program[0], string(exitErr.Stderr))
			}
			return nil, fmt.Errorf("command %q failed with no error message", program[0])
		} else {
			return nil, fmt.Errorf("failed to execute %q: %s", program[0], err)
		}
	}
	var decoded interface{}
	err = json.Unmarshal(resultJson, &decoded)
	if err != nil {
		return nil, fmt.Errorf("command %q produced invalid JSON: %s", program[0], err)
	}

	result := decoded.(map[string]interface{})
	return result, nil
}

func resourceExternalCreate(d *schema.ResourceData, meta interface{}) error {
	programI := d.Get("create").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	result, err := runCommand(programI, workingDir, query, d.Id())
	if err != nil {
		return fmt.Errorf("create: %s", err)
	}

	d.Set("result", result)
	d.SetId(result["id"].(string))
	log.Printf("[INFO] Created generic resource: %s", d.Id())

	return nil
}

func resourceExternalRead(d *schema.ResourceData, meta interface{}) error {
	programI := d.Get("read").([]interface{})
	log.Printf("[INFO] Number of command args [%d]", len(programI))
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	result, err := runCommand(programI, workingDir, query, d.Id())
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
	programI := d.Get("update").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	result, err := runCommand(programI, workingDir, query, d.Id())
	if err != nil {
		return fmt.Errorf("update: %s", err)
	}

	d.Set("result", result)
	log.Printf("[INFO] Updated resource: %s", d.Id())

	return nil
}
func resourceExternalDelete(d *schema.ResourceData, meta interface{}) error {
	programI := d.Get("delete").([]interface{})
	workingDir := d.Get("working_dir").(string)
	query := d.Get("query").(map[string]interface{})

	result, err := runCommand(programI, workingDir, query, d.Id())
	if err != nil {
		return fmt.Errorf("delete: %s", err)
	}

	log.Printf("[INFO] Deleted resource: %s", d.Id())
	d.Set("result", result)
	d.SetId("")

	return nil
}
