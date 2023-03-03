package service

import (
	"context"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/google/uuid"
	"github.com/ictsc/growi_client"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strings"
	"time"
)

type GrowiProblemSync struct {
	client                    *growi_client.GrowiClient
	path                      string
	authorId                  string
	problemWithInfoRepository repository.ProblemWithSyncTimeRepository
	problemRepository         repository.ProblemRepository
}

func NewGrowiProblemSyncService(
	client *growi_client.GrowiClient,
	path string,
	authorId string,
	problemWithInfoRepository repository.ProblemWithSyncTimeRepository,
	problemRepository repository.ProblemRepository,
) *GrowiProblemSync {
	return &GrowiProblemSync{
		client:                    client,
		path:                      path,
		authorId:                  authorId,
		problemWithInfoRepository: problemWithInfoRepository,
		problemRepository:         problemRepository,
	}
}

func (s *GrowiProblemSync) Sync(ctx context.Context) error {
	problems, err := s.problemRepository.GetAll()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get problems").Error())
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
			problemPage, err := s.client.GetPage(page.Path)
			if err != nil {
				log.Fatalf(errors.Wrapf(err, "Failed to get page").Error())
			}

			var matter = &entity.ProblemFrontMatter{}

			// frontmatter
			// TODO(k-shir0): フォーマットもチェックする
			_, err = frontmatter.Parse(strings.NewReader(problemPage.Revision.Body), matter)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(matter)
			//fmt.Println(string(body))

			// ここから先最終更新日と問題内容をキャッシュしておく
			newProblemWithInfo := &entity.ProblemWithSyncTime{
				Problem: entity.Problem{
					Code:              matter.Code,
					AuthorID:          uuid.MustParse(s.authorId),
					Title:             matter.Title,
					Body:              problemPage.Revision.Body,
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
					newProblemWithInfo.Problem.Base.CreatedAt = p.CreatedAt

					exists = true
					break
				}
			}

			// 既に存在すれば更新、無ければ作成し、失敗したならその問題はスキップする
			if exists {
				_, err = s.problemRepository.Update(&newProblemWithInfo.Problem)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				newProblemWithInfo.Problem.Base.CreatedAt = time.Now()

				_, err = s.problemRepository.Create(&newProblemWithInfo.Problem)
				if err != nil {
					log.Println(err)
					continue
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
