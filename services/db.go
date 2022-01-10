package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sakiib/clean-api/models"
)

const (
	MYSQL_PORT = "3306"
)

type DBConfig struct {
	Username     string
	Password     string
	PodName      string
	PodNamespace string
	DBName       string
	Port         string
	DBURL        string
}

func NewDBConfig(username, password, podName, podNamespace, dbName string) (DBInterface, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("username/password empty")
	}
	if len(podName) == 0 || len(podNamespace) == 0 {
		return nil, errors.New("mysql pod name/namespace empty")
	}
	if len(dbName) == 0 {
		return nil, errors.New("db name empty")
	}

	return &DBConfig{
		Username:     username,
		Password:     password,
		PodName:      podName,
		PodNamespace: podNamespace,
		DBName:       dbName,
		Port:         MYSQL_PORT,
		DBURL:        fmt.Sprintf("%s:%s@tcp(%s.%s.svc:%s)/%s", username, password, podName, podNamespace, MYSQL_PORT, dbName),
	}, nil
}

func (db *DBConfig) GetItems() ([]models.Product, error) {
	conn, err := sql.Open("mysql", db.DBURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	qry := "SELECT * FROM product"
	results, err := conn.Query(qry)
	if err != nil {
		return nil, err
	}

	var prodList []models.Product
	for results.Next() {
		var prod models.Product
		err = results.Scan(&prod.Id, &prod.Name, &prod.Price)
		if err != nil {
			return nil, err
		}
		prodList = append(prodList, prod)
	}

	return prodList, nil
}

func (db *DBConfig) DeleteItemWithID() error {
	conn, err := sql.Open("mysql", db.DBURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	qry := "DELETE FROM product WHERE id = 1"
	if _, err := conn.Exec(qry); err != nil {
		return err
	}

	return nil
}
