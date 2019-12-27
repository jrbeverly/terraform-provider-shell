package shell

type Config struct {
	WorkingDirectory string
	TempDirectory    string
	Variables        map[string]interface{}
}
