package serverchan

type Config struct {
	Enable bool   `json:"enable" yaml:"enable"`
	Key    string `yaml:"key" json:"key"`
}
