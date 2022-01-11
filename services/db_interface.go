package db

import "github.com/sakiib/vault-secrets-demo/models"

type DBInterface interface {
	GetItems() ([]models.Product, error)
	DeleteItemWithID() error
}
