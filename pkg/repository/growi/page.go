package growi

import (
	"encoding/json"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"net/http"
	"net/url"
)

type PageRepository struct {
	client *http.Client
	u      *url.URL
	token  string
}

func NewPageRepository(client *http.Client, u *url.URL, token string) *PageRepository {
	return &PageRepository{client: client, u: u, token: token}
}

type PageResponse struct {
	Page entity.Page `json:"page"`
}

func (r *PageRepository) Get(path string) (*entity.Page, error) {
	r.u.Path = "_api/v3/page"

	q := r.u.Query()
	q.Set("access_token", r.token)
	q.Set("path", path)

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

	// PageResponse に変換
	var pageResponse PageResponse
	err = json.NewDecoder(resp.Body).Decode(&pageResponse)
	if err != nil {
		return nil, err
	}

	return &pageResponse.Page, nil

}
