package pty

func spawnPtyShell(target, executor string, agent *util.AgentConfig) (int, int, error) {
	return agent.Pid, 1, errors.New("Not supported on JS agents")
}
