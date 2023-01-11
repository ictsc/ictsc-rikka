package repository

import "github.com/ictsc/ictsc-rikka/pkg/entity"

type SubordinatedPageRepository interface {
	GetAll() ([]entity.SubordinatedPage, error)
}
