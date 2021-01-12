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
	config.Agent.Sleep = 60

	return config
}
