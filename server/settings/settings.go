package settings

type ServerSettings struct {
	Port int
	Ip   string
}

func GetSettings() ServerSettings {
	return ServerSettings{
		5000,
		"localhost",
	}
}
