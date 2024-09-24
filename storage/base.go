package storage

import "github.com/knumor/qpoll/models"

// Storage is an interface for poll storage.
type Storage interface {
	Save(p models.Poll) error
	Load(id string) (models.Poll, error)
	LoadByCode(code string) (models.Poll, error)
}
