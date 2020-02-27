package settings

type ServerSettings struct {
	Port int
	Ip   string
}

func GetSettings() ServerSettings {
	return ServerSettings{
		Port: 5000,
		Ip:   "0.0.0.0",
	}
}
