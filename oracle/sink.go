package oracle

/*
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

func OracleSinkCreator(conf *domain.Config, database string) []string {
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	s := conf.SinkDb.User + ":" + conf.SinkDb.Password
	dsn := fmt.Sprintf("oracle://%s@%s:%s/%s", s, conf.SinkDb.Host, conf.SinkDb.Port, database)
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
	//db.Select(&schemas, "select USERNAME from sys.dba_users")
	db = db.Raw("select USERNAME from sys.dba_users").Scan(&schemas)
	for _, schema := range schemas {
		for _, tableRule := range conf.TableRule {
			matchResult, _ := regexp.MatchString(tableRule.Sink.Schema, schema)
			if matchResult {
				tableSql := fmt.Sprintf("select TABLE_NAME from dba_tables where owner = '%s'", schema)
				//db.Select(&tables, tableSql)
				db = db.Raw(tableSql).Scan(&tables)
				for _, table := range tables {
					wg.Add(1)
					var table = table
					go func() {
						matchResult1, _ := regexp.MatchString(tableRule.Sink.Table, table)
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

	return creatorSink(database, conf, db, dbT)

}

func creatorSink(database string, conf *domain.Config, db *gorm.DB, tables []string) []string {
	/*var wg sync.WaitGroup
	var m []string
	for _, t := range tables {
		wg.Add(1)
		t := t
		go func() {
			database, table, _ := strings.Cut(t, ".")
			sql := fmt.Sprintf("SELECT COLUMN_NAME as ColumnName,DATA_TYPE as DataType,COLUMN_TYPE as ColumnType,COLUMN_KEY as ColumnKey FROM information_schema.COLUMNS WHERE table_schema = '%s' AND table_name = '%s'", database, table)
			var t1 []domain.TableMessage
			db = db.Raw(sql).Scan(&t1)
			a := ""
			p := ""
			for _, message := range t1 {
				if message.ColumnKey == "PRI" {
					p = message.ColumnName
				}
				a = a + "`" + message.ColumnName + "` " + mymap[message.DataType] + ",\n"
			}
			a = a + "PRIMARY KEY(`" + p + "`) NOT ENFORCED\n"
			b := fmt.Sprintf("'connector' = 'jdbc',\n'driver' = 'com.mysql.jdbc.Driver',\n'url' = 'jdbc:mysql://%s:%s/%s?useUnicode=true&characterEncoding=utf8&zeroDateTimeBehavior=convertToNull&serverTimezone=GMT%%2B8',\n'username' = '%s',\n'password' = '%s',\n'table-name' = '%s'\n", conf.SinkDb.Host, conf.SinkDb.Port, database, conf.SinkDb.User, conf.SinkDb.Password, table)
			m = append(m, fmt.Sprintf("CREATE TABLE IF NOT EXISTS `default_catalog`.`%s`.`%s_sink`(\n%s) with (\n%s);\n", database, table, a, b))
			wg.Done()
		}()
		wg.Wait()
	}

	return m
	var wg sync.WaitGroup
	var m []string
	for _, t := range tables {
		wg.Add(1)
		t := t
		go func() {
			var pri string
			var t1 []domain.OracleTableMessage
			schema, table, _ := strings.Cut(t, ".")
			db = db.Raw(fmt.Sprintf("select   constraint_name   from   user_constraints   where   constraint_type='P' and  TABLE_name='%s'", table)).Scan(&pri)
			if pri == "" {
				panic(fmt.Sprintf("表'%s'不存在主键！", table))
			}
			sql := fmt.Sprintf("select column_name,data_type,data_length,DATA_PRECISION ,DATA_SCALE from user_tab_cols where table_name='%s'", table)
			db = db.Raw(sql).Scan(&t1)
			a := ""
			t := ""
			for _, message := range t1 {
				if message.DataType == "NUMBER" {
					t = numberTrans(message.ColumnName, message.DataLength, message.DataPrecision, message.DataScale)
				} else {
					t = oraclemap[message.DataType]
				}
				a = a + "`" + message.ColumnName + "` " + t + ",\n"
			}
			a = a + "PRIMARY KEY(`" + pri + "`) NOT ENFORCED\n"
			b := fmt.Sprintf("'connector' = 'jdbc',\n'driver' = 'oracle.jdbc.driver.OracleDriver',\n'url' = 'jdbc:oracle:thin:@%s:%s/%s',\n'username' = '%s',\n'password' = '%s',\n'table-name' = '%s'\n", conf.SinkDb.Host, conf.SinkDb.Port, database, conf.SinkDb.User, conf.SinkDb.Password, table)
			//b := fmt.Sprintf("'table-name' = '%s',\n'connector' = 'oracle-cdc',\n'hostname' = '%s',\n'port' = '%s',\n'username' = '%s',\n'password' = '%s',\n'database-name' = '%s'\n'schema-name' = '%s'", table, conf.SrcDb.Host, conf.SrcDb.Port, conf.SrcDb.User, conf.SrcDb.Password, database, schema)
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
*/
