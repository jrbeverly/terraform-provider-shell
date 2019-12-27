package shell

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"
	"path/filepath"

	"github.com/go-cmd/cmd"
	"github.com/matryer/try"

	"io/ioutil"
)

const RetryWaitTimeInSeconds = 30 * time.Second
const MaximumWaitTimeInSeconds = 5 * time.Minute

type CmdRunner struct {
	TemporaryDirectory string
	RetryMaximum       int
}

func (run *CmdRunner) convertToEnvVars(args map[string]interface{}, path string) []string {
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

func (run *CmdRunner) readDataFile(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	var decoded interface{}
	err = json.Unmarshal(byteValue, &decoded)
	if err != nil {
		return nil, fmt.Errorf("JSON file at %q produced invalid JSON: %s", path, err)
	}

	result := decoded.(map[string]interface{})
	return result, nil
}

func (run *CmdRunner) TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(run.TemporaryDirectory, prefix+hex.EncodeToString(randBytes)+suffix)
}

func (run *CmdRunner) runShellCommand(
	programI []interface{},
	workingDir string,
	query map[string]interface{},
	id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := try.Do(func(ampt int) (bool, error) {
		var err error
		result, err = run.runCmd(programI, workingDir, query, id, ampt)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, run.RetryMaximum, err)
			time.Sleep(RetryWaitTimeInSeconds)
		}
		return ampt < run.RetryMaximum, err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (run *CmdRunner) runCmd(
	programI []interface{},
	workingDir string,
	query map[string]interface{},
	id string,
	retry int) (map[string]interface{}, error) {
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

	data_file := run.TempFileName("shell-", ".tfjson")
	env := run.convertToEnvVars(query, data_file)
	env = append(env, fmt.Sprintf("TF_RETRY_COUNT=%d", retry))
	env = append(env, fmt.Sprintf("TF_RETRY=%t", retry > 0))
	env = append(env, fmt.Sprintf("TF_ID=%s", id))

	cmd := cmd.NewCmd(program[0], program[1:]...)
	cmd.Dir = workingDir
	cmd.Env = env

	statusChan := cmd.Start()

	go func() {
		<-time.After(MaximumWaitTimeInSeconds)
		cmd.Stop()
	}()

	status := <-statusChan
	err := status.Error
	if err != nil {
		stderr := strings.Join(status.Stderr, "\n")
		return nil, fmt.Errorf("Failed during execution %q: %s\n%s", program[0], err, stderr)
	}

	if _, err := os.Stat(data_file); os.IsNotExist(err) {
		stderr := strings.Join(status.Stderr, "\n")
		return nil, fmt.Errorf("Output from command was not recorded: %s\n%s", data_file, stderr)
	}

	if !status.Complete {
		return nil, fmt.Errorf("Timeout exception on %q", program[0])
	}

	log.Printf("[INFO] Data file name: %s", data_file)
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

	return run.readDataFile(data_file)
}
