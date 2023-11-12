package service

import (
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/ictsc/growi_client"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strings"
)

type GrowiNoticeSync struct {
	client           growi_client.Client
	path             string
	noticeRepository repository.NoticeRepository
}

func NewGrowiNoticeSyncService(
	client growi_client.Client,
	path string,
	noticeRepository repository.NoticeRepository,
) *GrowiNoticeSync {
	return &GrowiNoticeSync{
		client:           client,
		path:             path,
		noticeRepository: noticeRepository,
	}
}

func (s *GrowiNoticeSync) Sync() error {
	pages, err := s.client.GetSubordinatedPage(s.path)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to get subordinated list").Error())
	}

	notices, err := s.noticeRepository.GetAll()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get notices").Error())
	}

	paths, err := regexp.Compile(fmt.Sprintf(`^%s/`, s.path))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to compile regexp").Error())
	}

PageLoop:
	for _, page := range pages {
		split := strings.Split(page.Path, "/")
		end := split[len(split)-1]

		log.Println("syncing: ", page.Path)

		// _ で始まるページはスキップされる
		if strings.HasPrefix(end, "_") {
			log.Println("skipped because it starts with _")
			continue
		}
		if paths.MatchString(page.Path) {
			noticePage, err := s.client.GetPage(page.Path)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Failed to get notice page").Error())
			}

			var matter = &entity.NoticeFrontMatter{}
			body, err := frontmatter.Parse(strings.NewReader(noticePage.Revision.Body), matter)
			if err != nil {
				log.Fatal(err)
			}
			err = matter.Validate()
			if err != nil {
				fmt.Println(errors.Wrap(err, "Failed to validate frontmatter").Error())
				continue
			}

			newNotice := &entity.Notice{
				Base: entity.Base{
					UpdatedAt: noticePage.UpdatedAt,
					CreatedAt: noticePage.CreatedAt,
				}, Title: matter.Title,
				Body:     string(body),
				SourceId: end,
				Draft:    matter.Draft,
			}

			exist := false
			for _, p := range notices {
				if p.SourceId == newNotice.SourceId {
					if p.UpdatedAt.Equal(newNotice.UpdatedAt) {
						// 更新がないのでスキップ
						log.Println("skipped because it is not updated")
						continue PageLoop
					}

					newNotice.ID = p.ID
					exist = true
					break
				}
			}

			if !exist {
				// 新規追加処理
				if _, err := s.noticeRepository.Create(newNotice); err != nil {
					log.Fatal(err)
				}
				log.Println("created")
			} else {
				// 更新処理
				if _, err := s.noticeRepository.Update(newNotice, true); err != nil {
					log.Fatal(err)
				}
				log.Println("updated")
			}
		}
	}

	return nil
}
