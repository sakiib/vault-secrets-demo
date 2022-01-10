package db

import "github.com/sakiib/clean-api/models"

type DBInterface interface {
	GetItems() ([]models.Product, error)
	DeleteItemWithID() error
}
