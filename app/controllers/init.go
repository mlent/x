package controllers

import (
	"github.com/revel/revel"
	"gopkg.in/gorp.v1"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"app/app/models"
	"fmt"
	"strings"
)

func init() {
	revel.OnAppStart(InitDb)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).RollBack, revel.FINALLY)
}

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)

	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Could not find param:" + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func makeQueryString(args string) string {
	trimmedStr := strings.Trim(args,  " ");
	if trimmedStr != "" {
		return "?" + args
	}
	return args
}

func getConnectionString() string {
	host := getParamString("db.host", "localhost")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "root")
	pass := getParamString("db.password", "password")
	dbname := getParamString("db.name", "documents")
	protocol := getParamString("db.protocol", "tcp")

	preArgs := getParamString("dbargs", " ")
	dbargs := makeQueryString(preArgs)

	return fmt.Sprintf(
		"%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs,
	)
}

var InitDb func() = func() {
	connectionString := strings.Trim(getConnectionString(), " ")
	if db, err := sql.Open("mysql", connectionString); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		Dbm = &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
		}
	}
	
	defineDocumentTable(Dbm)
	if err := Dbm.CreateTablesIfNotExists(); err != nil {
		revel.ERROR.Fatal(err)
	}
}

func defineDocumentTable(dbm *gorp.DbMap) {
	table := dbm.AddTableWithName(models.Document{}, "documents").SetKeys(true, "id")
	table.ColMap("name").SetMaxSize(25)
}
