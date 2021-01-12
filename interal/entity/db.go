package entity

import (
	"fmt"
	"gin_demo/interal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

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
	m := interal.GetDbEnv()
	g.Dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		m["host"], m["user"], m["password"], m["dbname"], m["port"])
}


