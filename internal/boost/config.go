package boost

type TimeSpanConfig struct {
	Start string `json:"start" yaml:"start"`
	End   string `json:"end" yaml:"end"`
}

type Config struct {
	Enable          bool             `json:"enable" yaml:"enable"`
	CartInterval    string           `yaml:"cart_interval" json:"cart_interval" default:"2m"`
	ReserveInterval string           `yaml:"reserve_interval" json:"reserve_interval" default:"550ms"`
	RecheckInterval string           `yaml:"recheck_interval" json:"recheck_interval" default:"500ms"`
	ReorderInterval string           `yaml:"reorder_interval" json:"reorder_interval" default:"500ms"`
	WarmUpTimeSpan  []TimeSpanConfig `yaml:"warm_up_time_span" json:"warm_up_time_span"`
	TimeSpan        []TimeSpanConfig `yaml:"time_span" json:"time_span"`
}
