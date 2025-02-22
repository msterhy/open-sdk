package esx

type Config struct {
	Enable    bool     `yaml:"Enable"`
	Addresses []string `yaml:"Addresses"`
	Username  string   `yaml:"Username"`
	Password  string   `yaml:"Password"`
}
