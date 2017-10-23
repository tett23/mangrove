package mangrove_db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tett23/mangrove/assets"
	yaml "gopkg.in/yaml.v2"
)

var configFile = "config/mysql.yml"
var dsn string

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		panic("database is still initialized")
	}

	return db
}

func InitDatabase(envName string) (*gorm.DB, error) {
	dsn := getDSN(envName)
	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database(" + dsn + ")")
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	if envName != "production" {
		EnableLog(true)
	} else {
		EnableLog(false)
	}

	return db, err
}

func EnableLog(e bool) {
	db.LogMode(e)
}

type Logger interface {
	Print(v ...interface{})
}

func SetLogger(l Logger) {
	db.SetLogger(l)
}

func CloseConnection() {
	db.Close()
}

func getDSN(envName string) string {
	if dsn != "" {
		return dsn
	}

	bytes, err := assets.Asset(configFile)
	if err != nil {
		panic(err)
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}

	env := m[envName].(map[interface{}]interface{})

	address := env["address"]
	user := env["user"]
	password := env["password"]
	database := env["database"]
	protocol := env["protocol"]
	charset := env["charset"]

	dsn = fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=true", user, password, protocol, address, database, charset)

	return dsn
}
