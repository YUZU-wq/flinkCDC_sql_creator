package domain

type Config struct {
	SrcDb     SrcDb       `yaml:"srcDb"`
	SinkDb    SinkDb      `yaml:"sinkDb"`
	OutputDir interface{} `yaml:"outputDir"`
	TableRule []TableRule `yaml:"tableRule"`
	Config    []string    `yaml:"config"`
}

type SrcDb struct {
	Host     interface{} `yaml:"host"`
	Port     interface{} `yaml:"port"`
	User     interface{} `yaml:"user"`
	Password interface{} `yaml:"password"`
	Type     interface{} `yaml:"type"`
}

type SinkDb struct {
	User     interface{} `yaml:"user"`
	Password interface{} `yaml:"password"`
	Type     interface{} `yaml:"type"`
	Host     interface{} `yaml:"host"`
	Port     interface{} `yaml:"port"`
}

type TableRule struct {
	Database interface{} `yaml:"database"`
	Table    interface{} `yaml:"table"`
	Schema   interface{} `yaml:"schema"`
}
