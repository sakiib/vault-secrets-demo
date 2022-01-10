package main

import (
	"fmt"

	dbsvc "github.com/sakiib/clean-api/services"
	"k8s.io/klog"
)

func main() {
	fmt.Println("hello mysql")
	db, err := dbsvc.NewDBConfig("sakib", "12345", "mysql", "db", "test")
	if err != nil {
		panic(err)
	}

	items, err := db.GetItems()
	if err != nil {
		klog.Errorf("failed to get db items: %s", err.Error())
	}

	klog.Infoln("db items: ", items)

	if err := db.DeleteItemWithID(); err != nil {
		klog.Errorf("failed to delete item: %s", err.Error())
	}
}
