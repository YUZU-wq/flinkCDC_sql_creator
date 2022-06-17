package oracle

import (
	"flinkCDC_sql_creator/domain"
	"fmt"
	oracle "github.com/wdrabbit/gorm-oracle"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var oraclemap = make(map[string]string)

func InitMap() {
	oraclemap["DECIMAL"] = "DECIMAL"
	oraclemap["FLOAT"] = "FLOAT"
	oraclemap["BINARY_FLOAT"] = "FLOAT"
	oraclemap["DOUBLE PRECISION"] = "DOUBLE"
	oraclemap["BINARY_DOUBLE"] = "DOUBLE"
	oraclemap["DATE"] = "TIMESTAMP"
	oraclemap["CHAR"] = "STRING"
	oraclemap["NCHAR"] = "STRING"
	oraclemap["NVARCHAR2"] = "STRING"
	oraclemap["VARCHAR"] = "STRING"
	oraclemap["VARCHAR2"] = "STRING"
	oraclemap["CLOB"] = "STRING"
	oraclemap["NCLOB"] = "STRING"
	oraclemap["XMLType"] = "STRING"
	oraclemap["BLOB"] = "BYTES"
	oraclemap["ROWID"] = "BYTES"
	oraclemap["INTERVAL DAY TO SECOND"] = "BIGINT"
	oraclemap["INTERVAL YEAR TO MONTH"] = "BIGINT"
}

func OracleSrcCreator(conf *domain.Config, database string) []string {
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。

	s := "'" + conf.SrcDb.User + "/\"" + conf.SrcDb.Password + "\"'"
	dsn := fmt.Sprintf("oracle://%s@%s:%s/%s", s /*conf.SrcDb.User, conf.SrcDb.Password,*/, conf.SrcDb.Host, conf.SrcDb.Port, database)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。

	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 1 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		}),
		//SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	var schemas []string
	var tables []string
	var dbT []string
	var wg sync.WaitGroup
	//   查询 执行用Scan 和Find 一样
	db = db.Raw("select USERNAME from sys.dba_users").Scan(&schemas)
	for _, schema := range schemas {
		for _, tableRule := range conf.TableRule {
			matchResult, _ := regexp.MatchString(tableRule.Src.Schema, schema)
			if matchResult {
				tableSql := fmt.Sprintf("select TABLE_NAME from dba_tables where owner = '%s'", schema)
				db = db.Raw(tableSql).Scan(&tables)
				for _, table := range tables {
					wg.Add(1)
					var table = table
					go func() {
						matchResult1, _ := regexp.MatchString(tableRule.Src.Table, table)
						if matchResult1 {
							dbT = append(dbT, schema+"."+table)
						}
						wg.Done()
					}()
					wg.Wait()
				}
			}
		}
	}

	return creatorSrc(database, conf, db, dbT)

}

func creatorSrc(database string, conf *domain.Config, db *gorm.DB, tables []string) []string {
	var wg sync.WaitGroup
	var m []string
	for _, t := range tables {
		wg.Add(1)
		t := t
		go func() {
			var pri string
			var t1 []domain.OracleTableMessage
			schema, table, _ := strings.Cut(t, ".")
			db = db.Raw(fmt.Sprintf("select   constraint_name   from   user_constraints   where   constraint_type='P' and  TABLE_name=upper('%s');", table)).Scan(&pri)
			if pri == "" {
				panic(fmt.Sprintf("表'%s'不存在主键！", table))
			}
			sql := fmt.Sprintf("select column_name,data_type,data_length,DATA_PRECISION ,DATA_SCALE from user_tab_cols where table_name='%s'", table)
			db = db.Raw(sql).Scan(&t1)
			a := ""
			t := ""
			for _, message := range t1 {
				//todo 逻辑变更
				if message.DataType == "NUMBER" {
					t = numberTrans(message.ColumnName, message.DataLength, message.DataPrecision, message.DataScale)
				} else {
					t = oraclemap[message.DataType]
				}
				a = a + "`" + message.ColumnName + "` " + t + ",\n"
			}
			a = a + "PRIMARY KEY(`" + pri + "`) NOT ENFORCED\n"
			b := fmt.Sprintf("'table-name' = '%s',\n'connector' = 'oracle-cdc',\n'hostname' = '%s',\n'port' = '%s',\n'username' = '%s',\n'password' = '%s',\n'database-name' = '%s'\n'schema-name' = '%s'", table, conf.SrcDb.Host, conf.SrcDb.Port, conf.SrcDb.User, conf.SrcDb.Password, database, schema)
			for _, s := range conf.Config {
				b = b + ",\n" + s + "\n"
			}
			m = append(m, fmt.Sprintf("CREATE TABLE IF NOT EXISTS `default_catalog`.`%s`.`%s_src`(\n%s) with (\n%s);\n", schema, table, a, b))
			wg.Done()
		}()
		wg.Wait()
	}

	return m
}

func numberTrans(name string, l int, p int, s int) string {
	if l == 1 {
		return "BOOLEAN"
	}
	if p <= 0 && s <= 0 {
		if p-s < 3 {
			return "TINYINT"
		} else if p-s < 5 {
			return "SMALLINT"
		} else if p-s < 10 {
			return "INT"
		} else if p-s < 19 {
			return "BIGINT"
		} else if p-s >= 19 && p-s <= 38 {
			return "DECIMAL"
		} else {
			return "STRING"
		}
	} else if p > 0 && s > 0 {
		return "DECIMAL"
	} else {
		panic(name + "字段类型匹配失败！")
	}
}
