package model

// ProbeResult records isolated Moon CLI verification commands and outcome.
type ProbeResult struct {
	Module          string            `json:"module"`
	Version         string            `json:"version"`
	Target          string            `json:"target,omitempty"`
	ProbePath       string            `json:"probe_path"`
	MoonNew         CommandResult     `json:"moon_new"`
	MoonAdd         CommandResult     `json:"moon_add"`
	MoonCheck       CommandResult     `json:"moon_check"`
	MoonTest        CommandResult     `json:"moon_test"`
	MoonCheckTarget CommandResult     `json:"moon_check_target"`
	MoonBuildTarget CommandResult     `json:"moon_build_target"`
	Logs            map[string]string `json:"logs,omitempty"`
	Result          string            `json:"result"`
}

// CommandResult records one invoked command and its status.
type CommandResult struct {
	Command  []string `json:"command"`
	CWD      string   `json:"cwd"`
	ExitCode int      `json:"exit_code"`
	Status   State    `json:"status"`
	LogPath  string   `json:"log_path,omitempty"`
}
