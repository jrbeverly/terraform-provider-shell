package shell

import (
	"github.com/Brightspace/terraform-provider-shell/shell/api"
)

type Config struct {
	WorkingDirectory string
	Runner           api.CmdRunner
	Variables        map[string]interface{}
}
