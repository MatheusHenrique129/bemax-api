package mysql

import (
	"database/sql"
	"fmt"
	"time"

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

	// TODO validate if is necessary?
	dbClient.SetMaxIdleConns(25)
	dbClient.SetMaxIdleConns(5)
	dbClient.SetConnMaxLifetime(5 * time.Minute)
	dbClient.SetConnMaxIdleTime(1 * time.Minute)

	fmt.Println("successfully connected to MySQL database")
	return dbClient, nil
}

func createDB(driverName, user, pass, host, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, pass, host, name)
	return sql.Open(driverName, dsn)
}
