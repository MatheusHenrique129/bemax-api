package mysql

import (
	"database/sql"
	"fmt"

	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
	_ "github.com/go-sql-driver/mysql"
)

func NewMysql(config ports.MysqlConfig) (dbClient *sql.DB, err error) {
	dbClient, err = createDB(
		config.DriverName, config.UserName, config.UserPassword, config.HostName, config.DBName)
	if err != nil {
		_ = fmt.Errorf("could not create the database connection: %v", err)
		return nil, err
	}

	if err = dbClient.Ping(); err != nil {
		_ = dbClient.Close()
		_ = fmt.Errorf("could not create the database connection: %v", err)
		return nil, err
	}

	fmt.Println("successfully connected to MySQL database")
	return dbClient, nil
}

func createDB(driverName, user, pass, host, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, pass, host, name)
	return sql.Open(driverName, dsn)
}
