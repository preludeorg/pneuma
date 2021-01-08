package util

type Config struct {
	Agent struct {
		AESKey    string
		Range     string
		Contact   string
		Address   string
		Useragent string
		Sleep     int
		KillDate  string
	}
	Restraint struct {
		KillDate       string
		HttpKillSwitch string
		FileKillSwitch string
		AllowHost      []string
		AllowUser      []string
	}
}

func BuildConfig() *Config {
	config := &Config{}
	// Encryption
	config.Agent.AESKey = "abcdefghijklmnopqrstuvwxyz012345"

	// Contact
	config.Agent.Range = "red"
	config.Agent.Contact = "tcp"
	config.Agent.Address = "127.0.0.1:2323"
	config.Agent.Useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/87.0.4280.66 Safari/537.36"
	config.Agent.Sleep = 5

	// Restraints
	// Optionally remove a restraint by commenting it out or setting the value to nil.
	// Execution blocked after this date.
	config.Restraint.KillDate = "2077-01-01T00:00:00.000Z"
	// Execution blocked when the domain is resolved.
	//config.Restraint.HttpKillSwitch = "https://268d4c18eae4f534344dc2e6d7b1c72d6eacad3d34b43861e9b20d3e646f8798.com"
	// Execution blocked when file not found in user's home directory.
	//config.Restraint.FileKillSwitch = ".operator"
	// Execution blocked if the hostname is not listed here.
	//config.Restraint.AllowHost = []string{"hostname"}
	// Execution blocked if the agent's process owner is not listed here.
	//config.Restraint.AllowUser = []string{"user"}

	return config
}
