package growi

import (
	"encoding/json"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SubordinatedPageRepository struct {
	client *http.Client
	u      *url.URL
	path   string
	token  string
}

func NewSubordinatedPageRepository(client *http.Client, u *url.URL, path string, token string) *SubordinatedPageRepository {
	return &SubordinatedPageRepository{client: client, u: u, path: path, token: token}
}

type SubordinatedPagesResponse struct {
	SubordinatedPages []entity.SubordinatedPage `json:"subordinatedPages"`
}

func (r *SubordinatedPageRepository) GetAll() ([]entity.SubordinatedPage, error) {
	r.u.Path = "_api/v3/pages/subordinated-list"

	q := r.u.Query()
	q.Set("access_token", r.token)
	q.Set("path", r.path)

	r.u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", r.u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subordinatedPagesResponse SubordinatedPagesResponse
	err = json.Unmarshal(body, &subordinatedPagesResponse)
	if err != nil {
		return nil, err
	}

	return subordinatedPagesResponse.SubordinatedPages, nil
}
