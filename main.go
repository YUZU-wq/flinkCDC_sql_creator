package main

import (
	"bufio"
	"flag"
	"flinkCDC_sql_creator/domain"
	"flinkCDC_sql_creator/mysql"
	"flinkCDC_sql_creator/oracle"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
)

var configFile = flag.String("f", "conf/config.yaml", "the config file")

func main() {
	//从文件获取配置
	var wg sync.WaitGroup
	wg.Add(1)
	var c domain.Config
	yamlFile, err := ioutil.ReadFile(*configFile)
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Println("配置文件解析失败！")
		return
	}
	mysql.InitMap()
	oracle.InitMap()
	var src []string
	var sink []string
	go func() {
		defer wg.Done()
		switch c.SrcDb.Type {
		case "mysql":
			src = mysql.MysqlSrcCreator(&c)
		case "oracle":
			for _, rule := range c.TableRule {
				src = append(src, oracle.OracleSrcCreator(&c, rule.Src.Database)...)
			}
		default:
			fmt.Println("未匹配的数据库类型，等待支持。")
		}
	}()
	switch c.SinkDb.Type {
	case "mysql":
		sink = mysql.MysqlSinkCreator(&c)
	case "oracle":
		fmt.Println("flink-jdbc不支持oracle相关操作！")
	default:
		fmt.Println("未匹配的数据库类型，等待支持。")
	}

	//todo 进行类别判断 采用不同的生成方法
	/*mysql.InitMap()
	src := mysql.MysqlSrcCreator(&c)
	sink = mysql.MysqlSinkCreator(&c)*/
	filePath := c.OutputDir + "/create.sql"
	os.MkdirAll(c.OutputDir, os.ModeDir)
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file error=%v\n", err)
		return
	}
	defer file.Close()
	wg.Wait()
	bufWriter := bufio.NewWriter(file)

	limitChan := make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU())) // 最大并发协程数
	var mutex sync.Mutex

	for _, s1 := range src {
		limitChan <- struct{}{}
		wg.Add(1)

		s1 := s1
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("WriteData panic: %v,stack: %s\n", e, debug.Stack())
					// return
				}

				wg.Done()
				<-limitChan
			}()

			mutex.Lock() // 要加锁/解锁，否则 bufWriter.WriteString 写入数据有问题
			_, err := bufWriter.WriteString(s1)
			if err != nil {
				fmt.Printf("WriteDataToTxt WriteString err: %v\n", err)
				return
			}
			mutex.Unlock()
		}()
	}

	for _, s1 := range sink {
		limitChan <- struct{}{}
		wg.Add(1)

		s1 := s1
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("WriteData panic: %v,stack: %s\n", e, debug.Stack())
					// return
				}

				wg.Done()
				<-limitChan
			}()

			mutex.Lock() // 要加锁/解锁，否则 bufWriter.WriteString 写入数据有问题
			_, err := bufWriter.WriteString(s1)
			if err != nil {
				fmt.Printf("WriteData WriteString err: %v\n", err)
				return
			}
			mutex.Unlock()
		}()
	}
	wg.Wait()
	err = bufWriter.Flush()
	if err != nil {
		fmt.Printf("写文件失败: %v\n", err)
	}

}
