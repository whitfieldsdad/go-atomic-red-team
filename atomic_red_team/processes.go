package atomic_red_team

import (
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
)

type Process struct {
	Id          string     `json:"id"`
	Time        time.Time  `json:"time"`
	StartTime   *time.Time `json:"start_time"`
	User        *User      `json:"user"`
	PID         int        `json:"pid"`
	PPID        int        `json:"ppid"`
	Executable  *File      `json:"executable"`
	CommandLine string     `json:"command"`
	Argv        []string   `json:"argv"`
	ExitCode    *int       `json:"exit_code,omitempty"`
	Stdout      string     `json:"stdout,omitempty"`
	Stderr      string     `json:"stderr,omitempty"`
}

func GetProcess(pid int) (*Process, error) {
	p, err := sysinfo.Process(pid)
	if err != nil {
		return nil, err
	}
	info, err := p.Info()
	if err != nil {
		return nil, err
	}
	file, _ := GetFile(info.Exe)
	if err != nil {
		return nil, err
	}
	return &Process{
		Id:          NewUUID4(),
		Time:        time.Now(),
		StartTime:   &info.StartTime,
		PID:         info.PID,
		PPID:        info.PPID,
		Executable:  file,
		CommandLine: strings.Join(info.Args, " "),
		Argv:        info.Args,
	}, nil
}

func GetProcessAncestors(pid int) ([]Process, error) {
	var processes []Process
	process, err := GetProcess(pid)
	if err != nil {
		return nil, err
	}
	processes = append(processes, *process)

	pid = process.PPID
	for {
		process, err := GetProcess(pid)
		if err != nil {
			break
		}
		processes = append(processes, *process)
		if process.PPID == 0 {
			break
		}
		pid = process.PPID
	}
	return processes, nil
}

// IsElevated checks to see if the current process is either running with elevated privileges, or was started by an administrative user.
func IsElevated() (bool, error) {
	return isElevated()
}
