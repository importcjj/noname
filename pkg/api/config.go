package api

type Config struct {
	Cookie      string `yaml:"cookie" json:"cookie" required:"true"`
	Channel     string `yaml:"channel" json:"channel" default:"undefined"`
	APIVersion  string `yaml:"api_version" json:"api_version" default:"9.50.1"`
	APPVersion  string `yaml:"app_version" json:"app_version" default:"2.85.2"`
	ClientID    string `yaml:"client_id" json:"client_id" default:"3"`
	UserAgent   string `yaml:"ua" json:"ua" required:"true"`
	SID         string `yaml:"sid" json:"sid"  required:"true"`
	OpenID      string `yaml:"openid" json:"openid" required:"true"`
	DeviceID    string `yaml:"device_id" json:"device_id" required:"true"`
	DeviceToken string `yaml:"device_token" json:"device_token"  required:"true"`
	DebugTime   string `yaml:"debug_time" json:"debug_time"`
	Host        string `yaml:"host" json:"host" default:"maicai.api.ddxq.mobi"`
	Refer       string `yaml:"refer" json:"refer" default:"https://wx.m.ddxq.mobi/"`
	Origin      string `yaml:"origin" json:"origin" default:"https://wx.m.ddxq.mobi"`
}

func (c *Config) check() error {
	return nil
}
