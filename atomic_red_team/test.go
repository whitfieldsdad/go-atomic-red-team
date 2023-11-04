package atomic_red_team

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Test struct {
	Name                   string             `json:"name,omitempty" yaml:"name,omitempty"`
	AutoGeneratedGuid      string             `json:"auto_generated_guid,omitempty" yaml:"auto_generated_guid,omitempty"`
	Description            string             `json:"description,omitempty" yaml:"description,omitempty"`
	SupportedPlatforms     []string           `json:"supported_platforms,omitempty" yaml:"supported_platforms,omitempty"`
	InputArguments         map[string]ArgSpec `json:"input_arguments,omitempty" yaml:"input_arguments,omitempty"`
	DependencyExecutorName string             `json:"dependency_executor_name,omitempty" yaml:"dependency_executor_name,omitempty"`
	Dependencies           []Dependency       `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Executor               Executor           `json:"executor,omitempty" yaml:"executor,omitempty"`
	AttackTechniqueId      string             `json:"-" yaml:"-"`
	AttackTechniqueName    string             `json:"-" yaml:"-"`
}

func (t Test) DisplayName() string {
	return fmt.Sprintf("%s: %s - %s", t.AttackTechniqueId, t.AttackTechniqueName, t.Name)
}

func (t Test) IsManual() bool {
	return t.Executor.Name == "manual" || t.DependencyExecutorName == "manual"
}

func (t Test) MatchesAnyFilter(testFilters []TestFilter) bool {
	if len(testFilters) == 0 {
		return true
	}
	for _, testFilter := range testFilters {
		if testFilter.Matches(t) {
			return true
		}
	}
	return false
}

func (t Test) MatchesCurrentPlatform() bool {
	return t.MatchesPlatform(runtime.GOOS)
}

func (t Test) MatchesPlatform(platform string) bool {
	if platform == "darwin" {
		platform = "macos"
	}
	return slices.Contains(t.SupportedPlatforms, platform)
}

func (t Test) Run(ctx context.Context, atomicsDir string, opts *TestOptions) (*TestResult, error) {
	err := t.checkRequirements()
	if err != nil {
		return nil, err
	}

	var executedCommands []ExecutedCommand

	// Perform dependency resolution (TODO).
	if len(t.Dependencies) > 0 {
		log.Warnf("Test %s has %d dependencies, but dependency resolution is not supported yet", t.Name, len(t.Dependencies))
	}

	// Execute the primary test command.
	executedCommand, err := t.executeCommand(ctx, t.Executor.Command, t.Executor.Name, atomicsDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}
	executedCommands = append(executedCommands, *executedCommand)

	// Execute the cleanup command if one is specified.
	cleanupCommand := t.Executor.CleanupCommand
	if cleanupCommand != "" {
		executedCommand, err := t.executeCommand(ctx, cleanupCommand, t.Executor.Name, atomicsDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute cleanup command")
		}
		executedCommands = append(executedCommands, *executedCommand)
	}
	testResult := &TestResult{
		Id:               NewUUID4(),
		Time:             time.Now(),
		Test:             t,
		ExecutedCommands: executedCommands,
	}
	return testResult, nil
}

func (t Test) checkRequirements() error {
	executor := t.Executor
	if executor.Name == "manual" {
		return errors.New("manual tests are not supported")
	}
	if !t.MatchesCurrentPlatform() {
		return errors.New("unsupported platform")
	}
	if executor.ElevationRequired {
		elevated, err := IsElevated()
		if err != nil {
			return errors.Wrap(err, "failed to check if current process is elevated")
		}
		if !elevated {
			return errors.New("test requires elevation")
		}
	}
	return nil
}

func (t Test) executeCommand(ctx context.Context, command, commandType, atomicsDir string) (*ExecutedCommand, error) {
	kwargs := map[string]string{}
	for k, v := range t.InputArguments {
		kwargs[k] = v.DefaultValue
	}
	command = interpolateKwargs(command, kwargs)
	command, err := patchAtomicsDir(command, atomicsDir)
	if err != nil {
		log.Fatalf("Failed to patch PathToAtomicsDir: %s", err)
	}
	c, err := NewCommand(command, commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create command")
	}
	return c.Execute(ctx)
}

type ArgSpec struct {
	Description  string `json:"description" yaml:"description"`
	Type         string `json:"type" yaml:"type"`
	DefaultValue string `json:"default" yaml:"default"`
}

type Executor struct {
	Name              string `json:"name" yaml:"name"`
	ElevationRequired bool   `json:"elevation_required" yaml:"elevation_required"`
	Command           string `json:"command" yaml:"command"`
	CleanupCommand    string `json:"cleanup_command,omitempty" yaml:"cleanup_command,omitempty"`
}

type Dependency struct {
	Description      string `json:"description" yaml:"description"`
	PrereqCommand    string `json:"prereq_command" yaml:"prereq_command"`
	GetPrereqCommand string `json:"get_prereq_command" yaml:"get_prereq_command"`
	executorName     string `json:"-" yaml:"-"`
}

// TODO
func (d Dependency) Check(ctx context.Context) (*DependencyCheckResult, error) {
	panic("not implemented")
}

func (d Dependency) Resolve(ctx context.Context) (DependencyCheckResult, error) {
	panic("not implemented")
}

type DependencyCheckResult struct {
	Id               string            `json:"id" yaml:"id"`
	Time             time.Time         `json:"time" yaml:"time"`
	Dependency       Dependency        `json:"dependency" yaml:"dependency"`
	ExecutedCommands []ExecutedCommand `json:"executed_commands" yaml:"executed_commands"`
	Error            string            `json:"error" yaml:"error"`
	Resolved         bool              `json:"resolved" yaml:"resolved"`
}

func NewDependencyCheckResult(dependency Dependency, executedCommands []ExecutedCommand, resolved bool, err error) DependencyCheckResult {
	result := DependencyCheckResult{
		Id:               NewUUID4(),
		Time:             time.Now(),
		Dependency:       dependency,
		ExecutedCommands: executedCommands,
		Resolved:         resolved,
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

func patchAtomicsDir(command, atomicsDir string) (string, error) {
	if atomicsDir == "" {
		return "", errors.New("a directory is required")
	}
	command = strings.Replace(command, "PathToAtomicsFolder", atomicsDir, -1)
	return command, nil
}

func interpolateKwargs(command string, kwargs map[string]string) string {
	for k, v := range kwargs {
		k = fmt.Sprintf("#{%s}", k)
		command = strings.Replace(command, k, v, -1)
	}
	return command
}
