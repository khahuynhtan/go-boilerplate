package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name  string
	Email *string // pointer to allow null values
}

type Blog struct {
	gorm.Model
	Author  Author `gorm:"embedded"`
	UpVotes int32
}

type Token struct {
	gorm.Model
	Token string
}

var ModelsNeedToMigrate = []interface{}{
	&Author{},
	&Blog{},
	&Token{},
}
