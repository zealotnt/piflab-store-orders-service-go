package lib

import (
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}
