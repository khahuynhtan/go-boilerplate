package repositories

import (
	"myapp/models"
	"myapp/utils"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	// InsertCockroachData(in *entities.InsertCockroachDto) error
	GetList() ([]*models.Author, error)
	Create(*models.Author) error
}

type dbRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthorRepository() AuthorRepository {
	return &dbRepositoryImpl{
		db: utils.GetDB(),
	}
}

func (r *dbRepositoryImpl) GetList() ([]*models.Author, error) {
	var authors []*models.Author
	r.db.Find(&authors)
	return authors, nil
}

func (r *dbRepositoryImpl) Create(author *models.Author) error {
	return r.db.Create(author).Error
}
