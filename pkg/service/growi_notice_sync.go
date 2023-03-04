package service

import (
	"context"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/ictsc/growi_client"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strings"
	"time"
)

type GrowiNoticeSync struct {
	client                   *growi_client.GrowiClient
	path                     string
	noticeWithInfoRepository repository.NoticeWithSyncTimeRepository
	noticeRepository         repository.NoticeRepository
}

func NewGrowiNoticeSyncService(
	client *growi_client.GrowiClient,
	path string,
	noticeWithInfoRepository repository.NoticeWithSyncTimeRepository,
	noticeRepository repository.NoticeRepository,
) *GrowiNoticeSync {
	return &GrowiNoticeSync{
		client:                   client,
		path:                     path,
		noticeWithInfoRepository: noticeWithInfoRepository,
		noticeRepository:         noticeRepository,
	}
}

func (s *GrowiNoticeSync) Sync(ctx context.Context) error {
	notices, err := s.noticeRepository.GetAll()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get notices").Error())
	}

	pages, err := s.client.GetSubordinatedPage(s.path)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to get subordinated list").Error())
	}

	// ProblemPath 以下のやつだけ同期
	r := regexp.MustCompile(fmt.Sprintf(`^%s/`, s.path))

	for _, page := range pages {
		// _ で始まるパスを同期しないようしている
		split := strings.Split(page.Path, "/")
		end := split[len(split)-1]

		// どこのパスかどうかのログ
		fmt.Println(page.Path)

		// _ で始まるページは同期しない
		if strings.HasPrefix(end, "_") {
			fmt.Println("Sync Skip")
			continue
		}
		if r.MatchString(page.Path) {
			// redis キャッシュから取得し
			cachedNoticeWithInfo, err := s.noticeWithInfoRepository.Get(ctx, page.Path)
			if err != nil {
				fmt.Println("Not found in redis")
			} else {
				if cachedNoticeWithInfo.UpdatedAt == page.UpdatedAt {
					fmt.Println("Not updated")
					continue
				}
			}

			noticePage, err := s.client.GetPage(page.Path)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Failed to get notice page").Error())
			}

			var matter = &entity.NoticeFrontMatter{}

			body, err := frontmatter.Parse(strings.NewReader(noticePage.Revision.Body), matter)
			if err != nil {
				fmt.Println(errors.Wrap(err, "Failed to parse frontmatter").Error())
				continue
			}
			err = matter.Validate()
			if err != nil {
				fmt.Println(errors.Wrap(err, "Failed to validate frontmatter").Error())
				continue
			}

			fmt.Println(matter)
			fmt.Println(string(body))

			// ここから更新処理
			// 1. 最終日付を更新
			// 2. notice をキャッシュ
			newNoticeWithInfo := &entity.NoticeWithSyncTime{
				Notice: entity.Notice{
					Title:    matter.Title,
					Body:     string(body),
					SourceId: end,
					Draft:    matter.Draft,
				},
				UpdatedAt: page.UpdatedAt,
			}

			var exists = false
			for _, notice := range notices {
				// path を スラッシュで区切って一番最後の文字列
				// 例: /notice/2020/01/01/notice1 -> notice1
				if notice.SourceId == end {
					newNoticeWithInfo.ID = notice.ID
					newNoticeWithInfo.CreatedAt = notice.CreatedAt

					exists = true
					break

				}
			}

			// 既に存在する場合は更新
			if exists {
				_, err = s.noticeRepository.Update(&newNoticeWithInfo.Notice)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				newNoticeWithInfo.Notice.Base.CreatedAt = time.Now()

				_, err = s.noticeRepository.Create(&newNoticeWithInfo.Notice)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			err = s.noticeWithInfoRepository.Set(ctx, page.Path, *newNoticeWithInfo)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Failed to set notice to redis").Error())
			}
		}
	}

	return nil
}
