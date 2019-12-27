package shell

type Config struct {
	WorkingDirectory string
	Runner           CmdRunner
	Variables        map[string]interface{}
}
