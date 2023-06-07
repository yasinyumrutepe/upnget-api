package database

import (
	"auction/models"
	"auction/secret"
	"fmt"
	"log"
	"reflect"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type CstmTime time.Time

type Credentials struct {
	host    string
	port    string
	user    string
	pass    string
	sslMode string
	dBName  string
	schemas []string
}

type Database struct {
	myCredentials Credentials
	DB            *gorm.DB
}

var (
	Conn *Database
)

func (db *Database) Connect(dbname ...string) {
	if Conn != nil {
		dbCon, _ := db.DB.DB()
		dbCon.Close()
	}
	envPostEnv := secret.Env["myCredentials"].(map[string]interface{})
	db = &Database{}
	db.myCredentials = Credentials{
		host:    envPostEnv["host"].(string),
		port:    envPostEnv["port"].(string),
		user:    envPostEnv["user"].(string),
		pass:    envPostEnv["pass"].(string),
		sslMode: envPostEnv["sslMode"].(string),
		dBName:  envPostEnv["dbName"].(string),
	}
	if len(dbname) != 0 {
		db.myCredentials.dBName = dbname[0]
	}
	for _, v := range envPostEnv["schemas"].([]interface{}) {
		db.myCredentials.schemas = append(db.myCredentials.schemas, v.(string))
	}
	dsn := "host=" + db.myCredentials.host + " port=" + db.myCredentials.port + " user=" + db.myCredentials.user + " password=" + db.myCredentials.pass + " dbname=" + db.myCredentials.dBName + " sslmode=" + db.myCredentials.sslMode

	dba, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}
	db.DB = dba
	Conn = db
}

func (db *Database) Seed() error {
	var err error
	dbs := db.DB

	//DESC - All Seeds are created down in order
	models.User{}.Seed(dbs)
	models.Category{}.Seed(dbs)
	models.Brand{}.Seed(dbs)

	//DESC - All Seeds are called down in order

	return err
}

func (db *Database) Migrate() {
	var err error
	for {
		for _, schema := range db.myCredentials.schemas {

			Conn.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + " AUTHORIZATION " + Conn.myCredentials.user + ";")
			fmt.Println("Migrating schema: " + schema)
			for _, v := range migrateRelationList {
				if err := Conn.DB.SetupJoinTable(v.Model, v.Field, v.JoinTable); err != nil {
					fmt.Println(err.Error())
				}
			}
			err = Conn.DB.AutoMigrate(migrateModelList...)
			fmt.Println(err)

		}
		if err != nil {
			fmt.Println(err)
		}
		if err == nil {
			break
		}
	}
}

func (a Database) String() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s", a.myCredentials.user, a.myCredentials.pass, a.myCredentials.host, a.myCredentials.port, a.myCredentials.dBName)
}

func (db *Database) Close() {
	sqlDB, _ := db.DB.DB()
	sqlDB.Close()
}

func (db *Database) CreateSchema(schema string, user ...string) {
	var userName string
	if len(user) > 1 {
		panic("CreateSchema only accept one user")
	} else if len(user) == 0 {
		userName = db.myCredentials.user
	} else {
		userName = user[0]
	}
	if err := db.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + " AUTHORIZATION " + userName + ";").Error; err != nil {
		panic(err)
	}
}

func (db *Database) DropSchema(schema string) {
	if err := db.DB.Exec("DROP SCHEMA IF EXISTS " + schema + " CASCADE;").Error; err != nil {
		panic(err)
	}
}

func (db *Database) CreateTable(table interface{}) {
	if hasFlag := db.DB.Migrator().HasTable(table); !hasFlag {
		if err := db.DB.Migrator().CreateTable(table); err != nil {
			panic(err)
		}
	}
}

func (db *Database) DropTable(table interface{}) {
	if err := db.DB.Migrator().DropTable(table); err != nil {
		panic(err)
	}
}

func (db *Database) CreateJoinTable(model any, field string, joinTable any) {
	modelType := reflect.TypeOf(model).Elem().Kind()
	joinTableType := reflect.TypeOf(joinTable).Elem().Kind()
	if modelType != reflect.Struct && joinTableType != reflect.Struct {
		if err := db.DB.SetupJoinTable(model, field, joinTable); err != nil {
			panic(err)
		}
	}
}

func (db *Database) Add(structTemp interface{}) *gorm.DB {
	if tx := db.DB.Create(structTemp); tx.Error != nil {
		log.Printf("%s", tx.Error)
	}
	return db.DB
}

func (db *Database) Delete(structTemp interface{}, id uint) *gorm.DB {
	if tx := db.DB.Delete(structTemp, id); tx.Error != nil {
		log.Printf("%s", tx.Error)
	}
	return db.DB
}
