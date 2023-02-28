package service

import (
	"context"
	"github.com/ictsc/growi_client"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
)

type GrowiClientInit struct {
	client                       *growi_client.GrowiClient
	u                            *url.URL
	growiSessionCookieRepository repository.GrowiSessionCookieRepository
}

func NewGrowiClientInitService(
	client *growi_client.GrowiClient,
	u *url.URL,
	growiSessionCookieRepository repository.GrowiSessionCookieRepository,
) *GrowiClientInit {
	return &GrowiClientInit{
		client:                       client,
		u:                            u,
		growiSessionCookieRepository: growiSessionCookieRepository,
	}
}

var sessionCookieKey = "connect.sid"

func (s *GrowiClientInit) Init(ctx context.Context) error {
	// セッションクッキーを取得
	sessionCookieValue, err := s.growiSessionCookieRepository.Get(ctx)
	if err != nil {
		// 取得出来なかったときはログイン処理を行う
		err = s.client.Init()
		if err != nil {
			log.Fatal(errors.Wrap(err, "Failed to login").Error())
		}

		// connect.sid の Cookie を取得
		for _, cookie := range s.client.Jar.Cookies(s.u) {
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
	s.client.Jar.SetCookies(s.u, []*http.Cookie{
		{
			Name:  sessionCookieKey,
			Value: sessionCookieValue,
		},
	})

	return nil
}
