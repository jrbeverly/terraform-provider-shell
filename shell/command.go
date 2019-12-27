package shell

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-cmd/cmd"

	"io/ioutil"
)

const MaximumRetryWaitTimeInSeconds = 15 * time.Minute
const RetryWaitTimeInSeconds = 30 * time.Second
const MaximumWaitTimeInSeconds = 5 * time.Minute

/*
Commands should be done here:
	Underlying command for running
	Commands for each of the lifecycle

How do we handle this?
	Pass in an environment variable to a files
		1) Path to terraform state data
		2) Path to output state
	Provider directory for TMP dir
		This is where we will write the path to items
	This way we don't need to parse the outputs from the commands:
		'shell' vs 'external' vs 'cmd_exec'
*/

func writeInputState(path string, contents []byte) error {
	err := ioutil.WriteFile(path, contents, 0644)
	return err
}

func convertToEnvVars(args map[string]interface{}, path string) []string {
	i := 0
	vars := make([]string, len(args))
	for key, val := range args {
		vars[i] = fmt.Sprintf("%s=%s", key, val.(string))
		i++
	}
	vars = append(vars, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	vars = append(vars, fmt.Sprintf("TF_DATA_FILE=%s", path))
	return vars
}

func readDataFile(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	var decoded interface{}
	err = json.Unmarshal(byteValue, &decoded)
	if err != nil {
		return nil, fmt.Errorf("JSON file at %q produced invalid JSON: %s", path, err)
	}

	result := decoded.(map[string]interface{})
	return result, nil
}

func runCommand(programI []interface{}, workingDir string, tmpDir string, query map[string]interface{}, id string) (map[string]interface{}, error) {
	log.Printf("[INFO] Number of command args [%d]", len(programI))
	log.Printf("[INFO] Number of command env vars [%d]", len(query))
	program := make([]string, len(programI))
	for i, vI := range programI {
		log.Printf("[INFO] Program [%d]: %s", i, vI.(string))
		program[i] = vI.(string)
	}
	if len(program) == 0 {
		return nil, fmt.Errorf("No command has been provided")
	}

	path := fmt.Sprintf("%s/%s", tmpDir, "something.data.json")
	env := convertToEnvVars(query, path)

	cmd := cmd.NewCmd(program[0], program[1:]...)
	cmd.Dir = workingDir
	cmd.Env = append(env, fmt.Sprintf("TF_ID=%s", id))

	statusChan := cmd.Start()

	go func() {
		<-time.After(MaximumWaitTimeInSeconds)
		cmd.Stop()
	}()

	status := <-statusChan
	err := status.Error
	if err != nil {
		return nil, fmt.Errorf("Failed during execution %q: %s", program[0], err)
	}

	if !status.Complete {
		return nil, fmt.Errorf("Timeout exception on %q", program[0])
	}

	for _, out := range status.Stdout {
		log.Printf("[INFO] %s", out)
	}

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

	return readDataFile(path)
}
