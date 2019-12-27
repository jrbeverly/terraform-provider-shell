package shell

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
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

func runCommand(programI []interface{}, workingDir string, query map[string]interface{}, id string) (map[string]interface{}, error) {
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

	env := convertToEnvVars(query)

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

	resultJson := strings.Join(status.Stdout, " ")
	resultJson = strings.TrimSpace(resultJson)
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
	err = json.Unmarshal([]byte(resultJson), &decoded)
	if err != nil {
		return nil, fmt.Errorf("command %q produced invalid JSON: %s", program[0], err)
	}

	result := decoded.(map[string]interface{})
	return result, nil
}
