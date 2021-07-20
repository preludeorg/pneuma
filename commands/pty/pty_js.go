package pty

func spawnPtyShell(target, executor string, agent *util.AgentConfig) (int, int, error) {
	return agent.Pid, commands.ErrorExitStatus, errors.New("Not supported on JS agents")
}
