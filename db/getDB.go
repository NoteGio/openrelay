package db

import (
	"regexp"
	"github.com/notegio/openrelay/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Setup postgres dialect
	_ "github.com/jinzhu/gorm/dialects/mysql" // Setup mysql dialect
	"strings"
	"fmt"
)

func GetDB(connectionString, passwordURI string) (*gorm.DB, error) {
	password := common.GetSecret(passwordURI)
	connectionStringRegex := regexp.MustCompile("([^:]+)://([^@]+)@([^/]*)(/.*)?")
	match := connectionStringRegex.FindStringSubmatch(connectionString)
	if match == nil {
		return nil, fmt.Errorf("Parsing connection string '%v' failed", connectionString)
	}
	dbname := "OR_DEFAULT_DB"
	if cstringDbName := strings.TrimPrefix(match[4], "/"); cstringDbName != "" {
		dbname = cstringDbName
	}
	if match[1] == "postgres" {
		if dbname == "OR_DEFAULT_DB" {
			dbname = "postgres"
		}
		pgConnectionString := fmt.Sprintf(
			"host=%v dbname=%v sslmode=disable user=%v password=%v",
			match[3],
			dbname,
			match[2],
			password,
		)
		return gorm.Open("postgres", pgConnectionString)
	} else if match[1] == "mysql" {
		if dbname == "OR_DEFAULT_DB" {
			dbname = "mysql"
		}
		mysqlConnectionString := fmt.Sprintf(
			"%v:%v@%v/%v?charset=utf8&parseTime=True&loc=Local",
			match[2],
			password,
			match[3],
			dbname,
		)
		return gorm.Open("mysql", mysqlConnectionString)
	}
	return nil, fmt.Errorf("Unknown database type '%v'", match[1])
}
