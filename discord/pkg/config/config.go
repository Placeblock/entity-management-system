package config

type EMS struct {
	RestHost   string `yaml:"rest_host"`
	RestPort   string `yaml:"rest_port"`
	SocketHost string `yaml:"socket_host"`
	SocketPort string `yaml:"socket_port"`
}

type Config struct {
	Token           string `yaml:"token"`
	Guild           string `yaml:"guild"`
	Ems             EMS    `yaml:"ems"`
	TeamsCategoryID string `yaml:"teams-category-id"`
	EveryoneRoleID  string `yaml:"everyone-role-id"`
}
