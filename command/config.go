package command

type Server struct {
	ServerName string `toml:"server_name"`
	Hostname   string `toml:"hostname"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	TailFile   string `toml:"tail_file"`
}

type Config struct {
	TailFile string `toml:"tail_file"`
	Servers  map[string]Server `toml:"servers"`
}
