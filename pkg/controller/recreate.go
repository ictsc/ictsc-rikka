package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type RecreateController struct {
	problemSerice *service.ProblemService
	host          string
}

func NewRecreateController(problemService *service.ProblemService, host string) *RecreateController {
	return &RecreateController{
		problemSerice: problemService,
		host:          host,
	}
}

func (c *RecreateController) GetStatus(group *entity.UserGroup, probcode string) ([]byte, error) {
	_, err := c.problemSerice.FindByCode(probcode)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", c.host, group.TeamID, probcode)
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.New("backend Error")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("io Error")
	}
	return b, nil
}

func (c *RecreateController) CreateRequest(group *entity.UserGroup, probcode string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", c.host, group.TeamID, probcode)
	resp, err := http.Post(url, "", nil)

	if err != nil {
		return nil, errors.New("backend Error")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return nil, errors.New("Conflict")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("backend Error")
	}
	return c.GetStatus(group, probcode)
}
