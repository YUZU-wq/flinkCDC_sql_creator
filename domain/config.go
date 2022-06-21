package domain

import "net/url"

type Config struct {
	SrcDb     Db          `yaml:"srcDb"`
	SinkDb    Db          `yaml:"sinkDb"`
	OutputDir string      `yaml:"outputDir"`
	TableRule []TableRule `yaml:"tableRule"`
	Config    []string    `yaml:"config"`
}

type Db struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Type     string `yaml:"type"`
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

func (db Db) Trans() {
	db.User = url.QueryEscape(db.User)
	db.Password = url.QueryEscape(db.Password)
}
