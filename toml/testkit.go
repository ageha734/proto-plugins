package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

type Plugin struct {
	Name     string                 `toml:"name"`
	Platform map[string]interface{} `toml:"platform"`
}

type TestConfig struct {
	Name         string
	AfterInstall func(t *testing.T, shell *Shell) error
}

type TestResult struct {
	PluginName string
	Platform   string
	Supported  bool
	StartTime  time.Time
	EndTime    time.Time
	Success    bool
	Error      error
	LogFile    string
	Commands   []CommandLog
}

type CommandLog struct {
	Command   string
	Output    string
	Error     string
	StartTime time.Time
	EndTime   time.Time
	Success   bool
}

type Shell struct {
	t        *testing.T
	commands []CommandLog
}

func Run(config TestConfig) func(*testing.T) {
	return func(t *testing.T) {
		result := initializeTestResult(config.Name)
		var shell *Shell

		defer func() {
			finalizeTestResult(result, shell)
		}()

		printTestHeader(config.Name)

		plugin, tomlPathSource := loadPluginConfig(t, config.Name)
		platform := getPlatform()
		supportPlatforms := extractSupportedPlatforms(plugin)

		skip := !contains(supportPlatforms, platform)
		result.Platform = platform
		result.Supported = !skip

		printPlatformInfo(config.Name, supportPlatforms, platform, skip)

		if skip {
			t.Skipf("Platform %s not supported by plugin %s", platform, config.Name)
		}

		tempDir := createTempDirectory(t, config.Name)
		defer cleanupTempDirectory(tempDir)

		copyTomlFile(t, tomlPathSource, tempDir, config.Name)
		originalDir := changeToTempDirectory(t, tempDir)
		defer restoreOriginalDirectory(originalDir)

		shell = initializeShell(t)
		executePluginInstallation(shell, config.Name)
		executeAfterInstallTests(t, shell, config.AfterInstall)
	}
}

func initializeTestResult(pluginName string) *TestResult {
	return &TestResult{
		PluginName: pluginName,
		StartTime:  time.Now(),
	}
}

func finalizeTestResult(result *TestResult, shell *Shell) {
	result.EndTime = time.Now()
	result.Success = result.Error == nil
	if shell != nil {
		result.Commands = shell.commands
	}
	if !result.Success {
		result.LogFile = writeFailureLog(result)
	}
	printTestResult(result)
}

func loadPluginConfig(t *testing.T, pluginName string) (Plugin, string) {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		t.Fatal("Could not get caller information")
	}
	testDir := filepath.Dir(filename)

	tomlPathSource := filepath.Join(testDir, pluginName+".toml")
	if !filepath.IsAbs(tomlPathSource) {
		tomlPathSource = filepath.Clean(tomlPathSource)
	}

	content, err := os.ReadFile(tomlPathSource)
	if err != nil {
		t.Fatalf("Failed to read %s.toml: %v", pluginName, err)
	}

	var plugin Plugin
	if err := toml.Unmarshal(content, &plugin); err != nil {
		t.Fatalf("Failed to parse %s.toml: %v", pluginName, err)
	}

	return plugin, tomlPathSource
}

func extractSupportedPlatforms(plugin Plugin) []string {
	supportPlatforms := make([]string, 0, len(plugin.Platform))
	for platformName := range plugin.Platform {
		supportPlatforms = append(supportPlatforms, platformName)
	}
	return supportPlatforms
}

func createTempDirectory(t *testing.T, pluginName string) string {
	printStep("Creating temporary directory...")
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("proto-plugin-test-%s-", pluginName))
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tempDir
}

func cleanupTempDirectory(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("Failed to remove temp directory: %v", err)
	}
}

func copyTomlFile(t *testing.T, sourcePath, tempDir, pluginName string) {
	tomlPathDist := filepath.Join(tempDir, pluginName+".toml")
	if err := copyFile(sourcePath, tomlPathDist); err != nil {
		t.Fatalf("Failed to copy %s.toml: %v", pluginName, err)
	}
}

func changeToTempDirectory(t *testing.T, tempDir string) string {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	return originalDir
}

func restoreOriginalDirectory(originalDir string) {
	if err := os.Chdir(originalDir); err != nil {
		log.Printf("Failed to change back to original directory: %v", err)
	}
}

func initializeShell(t *testing.T) *Shell {
	return &Shell{t: t, commands: make([]CommandLog, 0)}
}

func executePluginInstallation(shell *Shell, pluginName string) {
	printStep("Setting up test environment...")
	shell.Exec("pwd")

	printStep("Adding plugin...")
	shell.Exec(fmt.Sprintf("proto plugin add %s source:./%s.toml", pluginName, pluginName))

	printStep("Installing plugin...")
	shell.Exec(fmt.Sprintf("proto install %s latest", pluginName))
}

func executeAfterInstallTests(t *testing.T, shell *Shell, afterInstall func(*testing.T, *Shell) error) {
	if afterInstall != nil {
		printStep("Running after-install tests...")
		if err := afterInstall(t, shell); err != nil {
			t.Fatalf("After install hook failed: %v", err)
		}
	}
}

func (s *Shell) Exec(command string) {
	s.t.Helper()
	printCommand(command)

	startTime := time.Now()
	cmd := exec.Command("sh", "-c", command)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()
	endTime := time.Now()

	commandLog := CommandLog{
		Command:   command,
		Output:    stdout.String(),
		Error:     stderr.String(),
		StartTime: startTime,
		EndTime:   endTime,
		Success:   err == nil,
	}
	s.commands = append(s.commands, commandLog)

	if err != nil {
		s.t.Fatalf("Command failed: %s, error: %v", command, err)
	}
}

func (s *Shell) ExecWithOutput(command string) (string, error) {
	s.t.Helper()
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	return string(output), err
}

func getPlatform() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	if !filepath.IsAbs(src) {
		src = filepath.Clean(src)
	}
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Printf("Failed to close source file: %v", err)
		}
	}()

	if !filepath.IsAbs(dst) {
		dst = filepath.Clean(dst)
	}
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			log.Printf("Failed to close destination file: %v", err)
		}
	}()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func printTestHeader(pluginName string) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("🧪 Testing Plugin: %s\n", pluginName)
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}

func printPlatformInfo(pluginName string, supportPlatforms []string, platform string, skip bool) {
	status := "✅ SUPPORTED"
	if skip {
		status = "❌ NOT SUPPORTED"
	}

	fmt.Printf("📋 Platform Information:\n")
	fmt.Printf("   Plugin: %s\n", pluginName)
	fmt.Printf("   Supported platforms: %v\n", supportPlatforms)
	fmt.Printf("   Current platform: %s\n", platform)
	fmt.Printf("   Status: %s\n", status)
	fmt.Printf("\n")
}

func printStep(step string) {
	fmt.Printf("🔄 %s\n", step)
}

func printCommand(command string) {
	fmt.Printf("   💻 Executing: %s\n", command)
}

func printTestResult(result *TestResult) {
	duration := result.EndTime.Sub(result.StartTime)
	status := "✅ PASSED"
	if !result.Success {
		status = "❌ FAILED"
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", 60))
	fmt.Printf("📊 Test Result Summary:\n")
	fmt.Printf("   Plugin: %s\n", result.PluginName)
	fmt.Printf("   Platform: %s\n", result.Platform)
	fmt.Printf("   Status: %s\n", status)
	fmt.Printf("   Duration: %v\n", duration.Round(time.Millisecond))
	if !result.Success && result.LogFile != "" {
		fmt.Printf("   📝 Failure log: %s\n", result.LogFile)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 60))
}

func writeFailureLog(result *TestResult) string {
	logDir := "test-logs"
	if err := os.MkdirAll(logDir, 0o750); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return ""
	}

	timestamp := result.EndTime.Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("%s_%s_%s_failure.log", result.PluginName, result.Platform, timestamp)
	logFilePath := filepath.Join(logDir, logFileName)

	if !filepath.IsAbs(logFilePath) {
		logFilePath = filepath.Clean(logFilePath)
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		log.Printf("Failed to create log file: %v", err)
		return ""
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()

	logContent := fmt.Sprintf(`# Test Failure Log
Plugin: %s
Platform: %s
Start Time: %s
End Time: %s
Duration: %v
Status: FAILED

## Test Details
Plugin Name: %s
Platform: %s
Supported: %t
Start Time: %s
End Time: %s
Duration: %v

## Command Execution Log
`,
		result.PluginName,
		result.Platform,
		result.StartTime.Format("2006-01-02 15:04:05"),
		result.EndTime.Format("2006-01-02 15:04:05"),
		result.EndTime.Sub(result.StartTime),
		result.PluginName,
		result.Platform,
		result.Supported,
		result.StartTime.Format("2006-01-02 15:04:05"),
		result.EndTime.Format("2006-01-02 15:04:05"),
		result.EndTime.Sub(result.StartTime),
	)

	for i, cmd := range result.Commands {
		logContent += fmt.Sprintf(`
### Command %d: %s
Start Time: %s
End Time: %s
Duration: %v
Success: %t

**Output:**
%s

**Error:**
%s

%s
`,
			i+1,
			cmd.Command,
			cmd.StartTime.Format("2006-01-02 15:04:05"),
			cmd.EndTime.Format("2006-01-02 15:04:05"),
			cmd.EndTime.Sub(cmd.StartTime),
			cmd.Success,
			cmd.Output,
			cmd.Error,
			strings.Repeat("-", 60),
		)
	}

	logContent += fmt.Sprintf(`
## Error Information
This test failed during execution. Please check the command execution log above for more details.

Generated at: %s
`, time.Now().Format("2006-01-02 15:04:05"))

	if _, err := file.WriteString(logContent); err != nil {
		log.Printf("Failed to write log content: %v", err)
		return ""
	}

	return logFilePath
}
