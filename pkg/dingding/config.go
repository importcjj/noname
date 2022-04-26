package dingding

type Config struct {
	Enable  bool   `json:"enable" yaml:"enable"`
	Keyword string `yaml:"keyword" json:"keyword"`
	Hook    string `json:"hook" yaml:"hook"`
}
