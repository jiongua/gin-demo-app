package entity

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type PostgreConfig struct {
	Host   string
	User   string
	Passwd string
	Dbname string
	Port   int
}

type Gorm struct {
	Dsn string
	once sync.Once
	db *gorm.DB
}

var gormDB = &Gorm{
	Dsn:  "",
	once: sync.Once{},
	db:   nil,
}

func init()  {
	MigrateDB(Db())
}

func Db() *gorm.DB {
	return gormDB.Db()
}

func (g *Gorm) Db() *gorm.DB {
	g.once.Do(g.connect)

	if g.db == nil {
		log.Fatal("entity: database not connected")
	}
	return g.db
}

func (g *Gorm) connect() {
	g.loadConfig()
	db, err := gorm.Open(postgres.Open(g.Dsn), &gorm.Config{
		Logger: GormLogger,
	})
	if err != nil {
		log.Fatalf("open postgre error: %s\n", err.Error())
	}
	log.Debugf("open postgre ok: %s\n", g.Dsn)
	g.db = db
	//MigrateDB(db)
}

func (g *Gorm) loadConfig() {
	dbConfig := PostgreConfig{}
	dbConfigName := ConfigPath()
	err := gonfig.GetConf(dbConfigName, &dbConfig)
	if err != nil {
		log.Fatalf("read postgre config error: %s\n", err.Error())
	}
	g.Dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.User, dbConfig.Passwd, dbConfig.Dbname, dbConfig.Port)
}

func ConfigPath() string {
	//rootDir, _ := filepath.Abs(filepath.Dir("."))
	//fmt.Printf("root_dir: %s", rootDir)
	//return path.Join(rootDir, "config", "pgsql.json")
	return "/Users/jiongua/zhihu/gin-demo-app/config/pgsql.json"
}

