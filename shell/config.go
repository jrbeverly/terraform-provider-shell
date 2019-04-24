package shell

type Config struct {
	WorkingDirectory string
	Variables        map[string]interface{}
	Prune            []string
}
