package mysql

import (
	"flinkCDC_sql_creator/domain"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"sync"
)

func MysqlSinkCreator(conf *domain.Config) map[string]string {
	timeout := "10s" //连接超时，10秒
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", conf.SinkDb.User, conf.SinkDb.Password, conf.SinkDb.Host, conf.SinkDb.Port, "information_schema", timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	var dataBases []string
	var tables []string
	var dbT []string
	var wg sync.WaitGroup
	//   查询 执行用Scan 和Find 一样
	db = db.Raw("show databases").Scan(&dataBases)
	for _, dataBase := range dataBases {
		for _, tableRule := range conf.TableRule {
			matchResult, _ := regexp.MatchString(tableRule.Sink.Database, dataBase)
			if matchResult {
				tableSql := fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.TABLES WHERE table_schema = '%s' AND table_type = 'BASE TABLE'", dataBase)
				db = db.Raw(tableSql).Scan(&tables)
				for _, table := range tables {
					wg.Add(1)
					var table = table
					go func() {
						matchResult1, _ := regexp.MatchString(tableRule.Sink.Table, table)
						if matchResult1 {
							dbT = append(dbT, dataBase+"."+table)
						}
						wg.Done()
					}()
					wg.Wait()
				}
			}
		}
	}

	return creatorSrc(conf, db, dbT)

}

func creatorSink(conf *domain.Config, db *gorm.DB, tables []string) map[string]string {
	var wg sync.WaitGroup
	m := make(map[string]string)
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
			m[t] = fmt.Sprintf("CREATE TABLE IF NOT EXISTS `default_catalog`.`%s`.`%s_sink`(\n%s) with (\n%s);\n", database, table, a, b)
			wg.Done()
		}()
		wg.Wait()
	}

	return m
}
