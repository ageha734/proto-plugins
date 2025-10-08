package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RunReport struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Node       ActionNode  `json:"node"`
	Operations []Operation `json:"operations"`
	Status     string      `json:"status"`
}

type ActionNode struct {
	Action string       `json:"action"`
	Params ActionParams `json:"params"`
}

type ActionParams struct {
	Target string `json:"target"`
}

type Operation struct {
	Meta OperationMeta `json:"meta"`
}

type OperationMeta struct {
	Type    string `json:"type"`
	Command string `json:"command,omitempty"`
}

type TargetIdentity struct {
	Task    string
	Project string
}

func coreDebug(message string) {
	fmt.Printf("::debug::%s\n", message)
}

func coreWarning(message string) {
	fmt.Printf("::warning::%s\n", message)
}

func coreStartGroup(title string) {
	fmt.Printf("::group::%s\n", title)
}

func coreEndGroup() {
	fmt.Println("::endgroup::")
}

func loadReport(workspaceRoot string) (*RunReport, error) {
	for _, fileName := range []string{"ciReport.json", "runReport.json"} {
		localPath := filepath.Join(".moon", "cache", fileName)
		reportPath := filepath.Join(workspaceRoot, localPath)

		coreDebug(fmt.Sprintf("Finding run report at %s", localPath))

		exists, err := fileExists(reportPath)
		if err != nil {
			return nil, err
		}

		if exists {
			coreDebug("Found!")
			data, err := os.ReadFile(reportPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read report file %s: %w", reportPath, err)
			}

			var report RunReport
			if err := json.Unmarshal(data, &report); err != nil {
				return nil, fmt.Errorf("failed to parse json report %s: %w", reportPath, err)
			}
			return &report, nil
		}
	}
	return nil, nil
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	report, err := loadReport(root)
	if err != nil {
		log.Fatalf("Failed to load run report: %v", err)
	}

	if report == nil {
		coreWarning("Run report does not exist, has `moon ci` or `moon run` ran?")
		return
	}

	for _, action := range report.Actions {
		if action.Node.Action != "run-task" {
			continue
		}

		targetIdentity := parseTarget(action.Node.Params.Target)
		target := fmt.Sprintf("%s:%s", targetIdentity.Project, targetIdentity.Task)

		command, _ := commandOf(action)
		stdout, stderr, err := readStatus(root, targetIdentity)
		if err != nil {
			log.Printf("Warning: could not read status for target %s: %v", target, err)
			continue
		}

		hasStdout := strings.TrimSpace(stdout) != ""
		hasStderr := strings.TrimSpace(stderr) != ""

		badge, ok := statusBadges[action.Status]
		if !ok {
			badge = action.Status // fallback
		}

		coreStartGroup(fmt.Sprintf("%s %s", badge, bold(target)))

		if command != "" {
			fmt.Println(blue(fmt.Sprintf("$ %s", command)))
		}

		if hasStdout {
			fmt.Println(stdBadges.out)
			fmt.Println(stdout)
		}

		if hasStderr {
			fmt.Println(stdBadges.err)
			fmt.Println(stderr)
		}

		coreEndGroup()
	}
}

func parseTarget(target string) TargetIdentity {
	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return TargetIdentity{Project: "unknown", Task: "unknown"}
	}
	return TargetIdentity{Project: parts[0], Task: parts[1]}
}

func commandOf(action Action) (string, bool) {
	for _, operation := range action.Operations {
		if operation.Meta.Type == "task-execution" {
			return operation.Meta.Command, true
		}
	}
	return "", false
}

func readStatus(workspaceRoot string, identity TargetIdentity) (string, string, error) {
	statusDir := filepath.Join(workspaceRoot, ".moon", "cache", "states", identity.Project, identity.Task)
	stdoutPath := filepath.Join(statusDir, "stdout.log")
	stderrPath := filepath.Join(statusDir, "stderr.log")

	readLogFile := func(path string) (string, error) {
		exists, err := fileExists(path)
		if err != nil || !exists {
			return "", err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	stdout, err := readLogFile(stdoutPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read stdout log: %w", err)
	}

	stderr, err := readLogFile(stderrPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read stderr log: %w", err)
	}

	return stdout, stderr, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var statusBadges = map[string]string{
	"running":              bgGreen(" RUNNING "),
	"passed":               bgGreen(" PASS "),
	"failed":               bgRed(" FAIL "),
	"timed-out":            bgRed(" TIMED OUT "),
	"aborted":              bgRed(" ABORTED "),
	"invalid":              bgRed(" INVALID "),
	"failed-and-abort":     bgRed(" FAILED AND ABORT "),
	"skipped":              bgBlue(" SKIP "),
	"cached":               bgBlue(" CACHED "),
	"cached-from-remote": bgBlue(" REMOTE CACHED "),
}

type stdBadgeType struct {
	out string
	err string
}

var stdBadges = stdBadgeType{
	out: bgDarkGray(fmt.Sprintf("　%s STDOUT　", green("⏺"))),
	err: bgDarkGray(fmt.Sprintf("　%s STDERR　", red("⏺"))),
}

func bgGreen(text string) string   { return fmt.Sprintf("\u001b[42m%s\u001b[49m", text) }
func bgRed(text string) string     { return fmt.Sprintf("\u001b[41m%s\u001b[49m", text) }
func bgBlue(text string) string    { return fmt.Sprintf("\u001b[44m%s\u001b[49m", text) }
func bgDarkGray(text string) string { return fmt.Sprintf("\u001b[48;5;236m%s\u001b[49m", text) }
func bold(text string) string      { return fmt.Sprintf("\u001b[1m%s\u001b[22m", text) }
func green(text string) string     { return fmt.Sprintf("\u001b[32m%s\u001b[39m", text) }
func red(text string) string       { return fmt.Sprintf("\u001b[31m%s\u001b[39m", text) }
func blue(text string) string      { return fmt.Sprintf("\u001b[34m%s\u001b[39m", text) }
