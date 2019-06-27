package command

type Server struct {
	ServerName     string `toml:"server_name"`
	Hostname       string `toml:"hostname"`
	Port           int    `toml:"port"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
	PrivateKeyPath string `toml:"private_key_path"`
	TailFile       string `toml:"tail_file"`
	TailFlags      string `toml:"tail_flags"`
}

type Config struct {
	TailFile  string            `toml:"tail_file"`
	Servers   map[string]Server `toml:"servers"`
	Slient    bool              `toml:"slient"`
	TailFlags string            `toml:"tail_flags"`
}
