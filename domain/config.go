package domain

type Config struct {
	SrcDb     SrcDb       `yaml:"srcDb"`
	SinkDb    SinkDb      `yaml:"sinkDb"`
	OutputDir string      `yaml:"outputDir"`
	TableRule []TableRule `yaml:"tableRule"`
	Config    []string    `yaml:"config"`
}

type SrcDb struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Type     string `yaml:"type"`
}

type SinkDb struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type TableRule struct {
	Src  Rule `yaml:"src"`
	Sink Rule `yaml:"sink"`
}

type TableMessage struct {
	ColumnName string
	DataType   string
	ColumnType string
	ColumnKey  string
}

type Rule struct {
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Schema   string `yaml:"schema"`
}

type OracleTableMessage struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	DataType      string `gorm:"column:DATA_TYPE"`
	DataLength    int    `gorm:"column:DATA_LENGTH"`
	DataPrecision int    `gorm:"column:DATA_PRECISION"`
	DataScale     int    `gorm:"column:DATA_SCALE"`
}

/*func (t TableMessage) create(dataBaseType string) string {
	switch dataBaseType {
	case "mysql":
		fmt.Println(1)
	case "oracle":
		fmt.Println(2)
	default:
		fmt.Println(0)
	}
}*/
