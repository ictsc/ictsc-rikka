package service

import (
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
)

type GrowiProblemSync struct {
	client            growi_client.Client
	path              string
	authorId          string
	problemRepository repository.ProblemRepository
}

func NewGrowiProblemSyncService(
	client growi_client.Client,
	path string,
	authorId string,
	problemRepository repository.ProblemRepository,
) *GrowiProblemSync {
	return &GrowiProblemSync{
		client:            client,
		path:              path,
		authorId:          authorId,
		problemRepository: problemRepository,
	}
}

func (s *GrowiProblemSync) Sync() error {
	pages, err := s.client.GetSubordinatedPage(s.path)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to get subordinated list").Error())
	}

	problems, err := s.problemRepository.GetAll()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get problems").Error())
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
			problemPage, err := s.client.GetPage(page.Path)
			if err != nil {
				log.Fatalf(errors.Wrapf(err, "Failed to get page").Error())
			}

			matter := &entity.ProblemFrontMatter{}
			_, err = frontmatter.Parse(strings.NewReader(problemPage.Revision.Body), matter)
			if err != nil {
				log.Fatal(err)
			}
			if err = matter.Validate(); err != nil {
				log.Fatal(err)
			}

			newProblem := &entity.Problem{
				Base: entity.Base{
					UpdatedAt: problemPage.UpdatedAt,
					CreatedAt: problemPage.CreatedAt,
				},
				Code:            matter.Code,
				AuthorID:        uuid.MustParse(s.authorId),
				Title:           matter.Title,
				Body:            problemPage.Revision.Body,
				Type:            matter.Type,
				Point:           matter.Point,
				SolvedCriterion: matter.SolvedCriterion,
			}

			if err := newProblem.DeleteMatterQuestionWithQuestionFieldAttach(); err != nil {
				log.Fatal(err)
			}

			if err := newProblem.Validate(); err != nil {
				log.Fatal(err)
			}

			exist := false
			for _, p := range problems {
				if p.Code == newProblem.Code {
					if p.UpdatedAt.Equal(newProblem.UpdatedAt) {
						// 更新がないのでスキップ
						log.Println("skipped because it is not updated")
						continue PageLoop
					}

					exist = true
					newProblem.ID = p.ID
					break
				}
			}

			if !exist {
				// 新規追加処理
				if _, err := s.problemRepository.Create(newProblem); err != nil {
					log.Fatal(err)
				}
				log.Println("created")
			} else {
				// 更新処理
				if _, err := s.problemRepository.Update(newProblem, true); err != nil {
					log.Fatal(err)
				}
				log.Println("updated")
			}
		}
	}

	return nil
}
