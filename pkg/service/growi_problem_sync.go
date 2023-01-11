package service

import (
	"context"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type GrowiProblemSync struct {
	client                       *http.Client
	u                            *url.URL
	path                         string
	username                     string
	password                     string
	token                        string
	growiSessionCookieRepository repository.GrowiSessionCookieRepository
	problemWithInfoRepository    repository.ProblemWithInfoRepository
	pageRepository               repository.PageRepository
	subordinatedPageRepository   repository.SubordinatedPageRepository
	problemRepository            repository.ProblemRepository
}

func NewGrowiProblemSyncService(
	client *http.Client,
	u *url.URL,
	path string,
	username string,
	password string,
	token string,
	growiSessionCookieRepository repository.GrowiSessionCookieRepository,
	problemWithInfoRepository repository.ProblemWithInfoRepository,
	pageRepository repository.PageRepository,
	subordinatedPageRepository repository.SubordinatedPageRepository,
	problemRepository repository.ProblemRepository,
) *GrowiProblemSync {
	return &GrowiProblemSync{
		client:                       client,
		u:                            u,
		path:                         path,
		username:                     username,
		password:                     password,
		token:                        token,
		growiSessionCookieRepository: growiSessionCookieRepository,
		problemWithInfoRepository:    problemWithInfoRepository,
		pageRepository:               pageRepository,
		subordinatedPageRepository:   subordinatedPageRepository,
		problemRepository:            problemRepository,
	}
}

var sessionCookieKey = "connect.sid"

func (s *GrowiProblemSync) Init(ctx context.Context) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	s.client.Jar = jar

	// セッションクッキーを取得
	sessionCookieValue, err := s.growiSessionCookieRepository.Get(ctx)
	if err != nil {
		// 取得出来なかったときはログイン処理を行う
		csrfToken, err := getCsrfToken(*s.u, *s.client)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Failed to get csrf token").Error())
		}

		err = doLogin(*s.u, *s.client, s.username, s.password, csrfToken)
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "Failed to login").Error())
		}

		// connect.sid の Cookie を取得
		for _, cookie := range jar.Cookies(s.u) {
			if cookie.Name == "connect.sid" {
				err = s.growiSessionCookieRepository.Set(ctx, cookie.Value)
				if err != nil {
					log.Fatal(errors.Wrap(err, "Failed to set session cookie").Error())
				}
			}
		}

		return err
	}

	// sessionCookieValue を connect.sid の Cookie にセット
	jar.SetCookies(s.u, []*http.Cookie{
		{
			Name:  sessionCookieKey,
			Value: sessionCookieValue,
		},
	})

	return nil
}

func (s *GrowiProblemSync) Sync(ctx context.Context) error {
	problems, err := s.problemRepository.GetAll()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get problems").Error())
	}

	pages, err := s.subordinatedPageRepository.GetAll()
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to get subordinated list").Error())
	}

	// Path 以下のやつだけ同期
	r := regexp.MustCompile(fmt.Sprintf(`^%s/`, s.path))

	for _, page := range pages {
		if r.MatchString(page.Path) {
			// redis キャッシュから取得し
			cachedProblemWithInfo, err := s.problemWithInfoRepository.Get(ctx, page.Path)
			if err != nil {
				fmt.Println("Not found in redis")
			} else {
				// 更新日付が一緒なら同期しない
				if cachedProblemWithInfo.UpdatedAt == page.UpdatedAt {
					fmt.Println("Not updated")
					continue
				}
			}

			fmt.Println("Updating...")
			fmt.Println(page.Path)

			// 個別ページを取得
			problemPage, err := s.pageRepository.Get(page.Path)
			if err != nil {
				log.Fatalf(errors.Wrapf(err, "Failed to get page").Error())
			}

			var matter = &entity.Matter{}

			// frontmatter
			// TODO(k-shir0): フォーマットもチェックする
			body, err := frontmatter.Parse(strings.NewReader(problemPage.Revision.Body), matter)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(matter)
			fmt.Println(string(body))

			// ここから先最終更新日と問題内容をキャッシュしておく
			newProblemWithInfo := &entity.ProblemWithInfo{
				Problem: entity.Problem{
					Code:              matter.Code,
					AuthorID:          uuid.MustParse(matter.AuthorId),
					Title:             matter.Title,
					Body:              string(body),
					Point:             matter.Point,
					PreviousProblemID: nil,
					SolvedCriterion:   matter.SolvedCriterion,
				},
				UpdatedAt: page.UpdatedAt,
			}

			var exists = false
			for _, p := range problems {
				if p.Code == newProblemWithInfo.Problem.Code {
					newProblemWithInfo.Problem.ID = p.ID

					exists = true
					break
				}
			}

			// 既に存在すれば更新、無ければ作成
			if exists {
				_, err = s.problemRepository.Update(&newProblemWithInfo.Problem)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				_, err = s.problemRepository.Create(&newProblemWithInfo.Problem)
				if err != nil {
					log.Fatal(err)
				}
			}

			// キャッシュを行う
			err = s.problemWithInfoRepository.Set(ctx, page.Path, *newProblemWithInfo)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func doLogin(u url.URL, client http.Client, username string, password string, csrfToken string) error {
	u.Path = "/login"

	// body form-urlencoded
	form := url.Values{}
	form.Add("loginForm[username]", username)
	form.Add("loginForm[password]", password)
	form.Add("_csrf", csrfToken)

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getCsrfToken(u url.URL, client http.Client) (string, error) {
	u.Path = "/login"

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	node, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// html > body attr の中から csrfToken を取得
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "html" {
			for d := c.FirstChild; d != nil; d = d.NextSibling {
				if d.Data == "body" {
					for _, attr := range d.Attr {
						if attr.Key == "data-csrftoken" {
							return attr.Val, nil
						}
					}
				}
			}
		}
	}

	return "", nil
}
