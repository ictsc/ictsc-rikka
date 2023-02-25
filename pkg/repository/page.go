package repository

import "github.com/ictsc/ictsc-rikka/pkg/entity"

type PageRepository interface {
	Get(path string) (*entity.Page, error)
}
