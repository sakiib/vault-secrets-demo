package main

import (
	"flag"
	"os"
	"path/filepath"

	dbsvc "github.com/sakiib/vault-secrets-demo/services"
	"gopkg.in/fsnotify.v1"
	"k8s.io/klog"
)

var (
	username string
	password string
	watcher  *fsnotify.Watcher

	user        = flag.String("app-username", "sql-user", "app secret username")
	pass        = flag.String("app-password", "sql-pass", "app secret password")
	secretsPath = flag.String("secrets-path", "/secrets-store/sql-creds", "fs path")
)

func main() {
	flag.Parse()

	if len(*user) == 0 || len(*pass) == 0 {
		klog.Fatal("secret name not provided")
	}

	userFile := filepath.Join(*secretsPath, *user)
	passFile := filepath.Join(*secretsPath, *pass)

	user, err := getSecret(userFile)
	if err != nil {
		klog.Fatalf("failed to get username file: %v", err)
	}

	pass, err := getSecret(passFile)
	if err != nil {
		klog.Fatalf("failed to get password file: %v", err)
	}

	username = user
	password = pass

	db, err := dbsvc.NewDBConfig(username, password, "sql", "db", "test")
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

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err = watcher.Add(*secretsPath); err != nil {
		klog.Fatalf("failed to add secret path to file watcher: %v", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-watcher.Events:
				user, err := getSecret(userFile)
				if err != nil {
					klog.Fatalf("failed to get username file: %v", err)
				}
				if username != user {
					username = user
					klog.Infoln(username)
				}

				pass, err := getSecret(passFile)
				if err != nil {
					klog.Fatalf("failed to get password file: %v", err)
				}
				if password != pass {
					password = pass
					klog.Infoln("in: ", password)
				}

			case err := <-watcher.Errors:
				klog.Error("failed to watch: %v", err)
			}
		}
	}()

	<-done
}

func getSecret(secretFile string) (string, error) {
	if _, err := os.Stat(secretFile); err != nil {
		return "", err
	}

	secret, err := os.ReadFile(secretFile)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
