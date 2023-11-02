package atomic_red_team

import (
	"runtime"

	"github.com/pkg/errors"
)

var (
	WindowsPowerShell = "powershell"
	PowerShellCore    = "pwsh"
	PowerShell        = getPowerShellCommandType()
	CommandPrompt     = "command_prompt"
	Sh                = "sh"
	Bash              = "bash"
)

var (
	DefaultCommandType = getDefaultCommandType()
)

func getDefaultCommandType() string {
	if runtime.GOOS == "windows" {
		return CommandPrompt
	}
	return Bash
}

func getPowerShellCommandType() string {
	if runtime.GOOS == "windows" {
		return WindowsPowerShell
	}
	return PowerShellCore
}

var (
	commandTypes = []string{PowerShell, CommandPrompt, Sh, Bash}
)

var (
	commandShims = map[string][]string{
		WindowsPowerShell: []string{"powershell", "-ExecutionPolicy", "Bypass", "-Command"},
		PowerShellCore:    []string{"pwsh", "-Command"},
		CommandPrompt:     []string{"cmd", "/c"},
		Sh:                []string{"sh", "-c"},
		Bash:              []string{"bash", "-c"},
	}
)

func prepareCommand(command, commandType string) ([]string, error) {
	commandType, err := parseCommandType(commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command type")
	}
	shim, ok := commandShims[commandType]
	if !ok {
		return nil, errors.Errorf("invalid command type: %s", commandType)
	}
	return append(shim, command), nil
}

func parseCommandType(commandType string) (string, error) {
	for _, validCommandType := range commandTypes {
		if commandType == validCommandType {
			return commandType, nil
		}
	}
	return "", errors.Errorf("invalid command type: %s", commandType)
}
