package atomic_red_team

import (
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/exp/slices"
)

type ProcessFilter struct {
	PIDs             []int    `json:"pids,omitempty"`
	PPIDs            []int    `json:"ppids,omitempty"`
	Names            []string `json:"names,omitempty"`
	ExecutablePaths  []string `json:"executable_paths,omitempty"`
	ExecutableNames  []string `json:"executable_names,omitempty"`
	ExecutableHashes []string `json:"executable_hashes,omitempty"`
}

func (f ProcessFilter) Matches(p Process) bool {
	if len(f.PIDs) > 0 && !slices.Contains(f.PIDs, p.PID) {
		return false
	}
	if len(f.PPIDs) > 0 && !slices.Contains(f.PPIDs, p.PPID) {
		return false
	}
	if len(f.Names) > 0 && !slices.Contains(f.Names, p.Name) {
		return false
	}
	if len(f.ExecutablePaths) > 0 && !slices.Contains(f.ExecutablePaths, p.File.Path) {
		return false
	}
	if len(f.ExecutableNames) > 0 && !slices.Contains(f.ExecutableNames, p.File.Name) {
		return false
	}
	if len(f.ExecutableHashes) > 0 {
		if p.File == nil || p.File.Hashes == nil {
			return false
		}
		hashes := p.File.Hashes.List()
		for _, a := range hashes {
			for _, b := range f.ExecutableHashes {
				if strings.EqualFold(a, b) {
					return true
				}
			}
		}
		return false
	}
	return true
}

type Process struct {
	Id          string     `json:"id"`
	Time        time.Time  `json:"time"`
	Name        string     `json:"name,omitempty"`
	PID         int        `json:"pid"`
	PPID        int        `json:"ppid"`
	CWD         string     `json:"cwd,omitempty"`
	File        *File      `json:"file,omitempty"`
	CommandLine string     `json:"command_line,omitempty"`
	Argv        []string   `json:"argv,omitempty"`
	User        string     `json:"user,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	ExitCode    *int       `json:"exit_code,omitempty"`
	Stdout      string     `json:"stdout,omitempty"`
	Stderr      string     `json:"stderr,omitempty"`
}

func calculateProcessUUID(pid, ppid int) string {
	factors := map[string]interface{}{
		"pid":  pid,
		"ppid": ppid,
	}
	return NewUUID5(factors)
}

func (p Process) GetArtifactType() ArtifactType {
	return ProcessArtifactType
}

func (p Process) ToProcessEvents() []LightweightProcessEvent {
	var events []LightweightProcessEvent
	if p.StartTime != nil {
		events = append(events, LightweightProcessEvent{
			Event: Event{
				Id:         NewUUID4(),
				Time:       *p.StartTime,
				ObjectType: "process",
				EventType:  string(ProcessEventTypeStart),
			},
			PID:  p.PID,
			PPID: p.PPID,
		})
	}
	if p.ExitTime != nil {
		events = append(events, LightweightProcessEvent{
			Event: Event{
				Id:         NewUUID4(),
				Time:       *p.ExitTime,
				ObjectType: "process",
				EventType:  string(ProcessEventTypeExit),
			},
			PID:  p.PID,
			PPID: p.PPID,
		})
	}
	return events
}

func GetProcessFamily(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetProcessFamilyPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process family PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetProcessFamilyPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetFamilyMembers(pid)
}

func GetParentProcess(pid int) (*Process, error) {
	ppid, err := GetPPID(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get parent process")
	}
	return GetProcess(ppid)
}

func GetPPID(pid int) (int, error) {
	p, err := ps.FindProcess(pid)
	if err != nil {
		return -1, errors.Wrap(err, "failed to get process")
	}
	return p.PPid(), nil
}

func GetProcessAncestors(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetProcessAncestorPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ancestor PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetProcessAncestorPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetAncestors(pid)
}

func GetProcessDescendants(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetProcessDescendantPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get descendant PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetProcessDescendantPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetDescendants(pid)
}

func GetChildProcesses(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetChildProcessPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get child process PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetChildProcessPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetChildren(pid)
}

func GetProcessSiblings(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetProcessSiblingPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sibling PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetProcessSiblingPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetSiblings(pid)
}

func GetProcessSiblingDescendants(pid int, filter *ProcessFilter) ([]Process, error) {
	pids, err := GetProcessSiblingDescendantPIDs(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sibling descendant PIDs")
	}
	return GetProcesses(pids, filter)
}

func GetProcessSiblingDescendantPIDs(pid int) ([]int, error) {
	tree, err := GetLightweightProcessTree()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process tree")
	}
	return tree.GetSiblingDescendants(pid)
}

func GetPIDs() ([]int, error) {
	pids, err := ps.Processes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list processes")
	}
	var results []int
	for _, p := range pids {
		results = append(results, p.Pid())
	}
	return results, nil
}

func GetProcesses(pids []int, filter *ProcessFilter) ([]Process, error) {
	if pids == nil {
		var err error
		pids, err = GetPIDs()
		if err != nil {
			return nil, err
		}
	}
	var wg sync.WaitGroup
	processes := make(chan Process, len(pids))
	for _, pid := range pids {
		wg.Add(1)
		go doGetProcess(pid, processes, &wg)
	}
	wg.Wait()
	close(processes)

	var results []Process
	for p := range processes {
		if filter != nil && !filter.Matches(p) {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}

func doGetProcess(pid int, processes chan<- Process, wg *sync.WaitGroup) {
	defer wg.Done()
	p, err := GetProcess(pid)
	if err != nil {
		log.Debugf("Failed to collect process metadata (PID: %d): %s", pid, err)
		return
	}
	processes <- *p
}

func GetProcess(pid int) (*Process, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	var (
		startTimePtr *time.Time
		cmd          string
		file         *File
	)
	ppid, _ := p.Ppid()
	name, _ := p.Name()
	argv, _ := p.CmdlineSlice()
	if argv != nil {
		cmd = strings.Join(argv, " ")
	}
	cwd, _ := p.Cwd()
	path, err := p.Exe()
	if err == nil {
		file, _ = GetFile(path, nil)
	}
	startTimeMs, err := p.CreateTime()
	if err == nil {
		startTime := time.Unix(0, startTimeMs*int64(time.Millisecond))
		startTimePtr = &startTime
	}
	o := &Process{
		Id:          calculateProcessUUID(pid, int(ppid)),
		Time:        time.Now(),
		Name:        name,
		PID:         pid,
		PPID:        int(ppid),
		CWD:         cwd,
		CommandLine: cmd,
		Argv:        argv,
		File:        file,
		User:        CurrentUser,
		StartTime:   startTimePtr,
	}
	return o, nil
}

func CurrentProcessIsElevated() (bool, error) {
	log.Infof("Checking if current process is elevated...")
	elevated, err := currentProcessIsElevated()
	if err != nil {
		return false, errors.Wrap(err, "failed to check if current process is elevated")
	}
	log.Infof("Current process is elevated: %t", elevated)
	return elevated, nil
}
