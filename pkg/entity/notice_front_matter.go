package entity

import "errors"

type NoticeFrontMatter struct {
	Title string `yaml:"title"`
	Draft bool   `yaml:"draft"`
}

func (n *NoticeFrontMatter) Validate() error {
	if n.Title == "" {
		return errors.New("title must not be empty")
	}

	return nil
}
