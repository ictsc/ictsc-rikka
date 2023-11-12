package service

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ictsc/growi_client"
	gentity "github.com/ictsc/growi_client/entity"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mock"
	"testing"
	"time"
)

func TestGrowiProblemSync_Sync(t *testing.T) {
	tests := []struct {
		name     string
		mockInit func(growiClient *mockGrowiClient, problemRepo mock.MockProblemRepository)
		wantErr  bool
	}{
		{
			name: "ページ名が _ で始まるページがスキップされる",
			mockInit: func(growiClient *mockGrowiClient, problemRepo mock.MockProblemRepository) {
				growiClient.getSubordinatedPageReturn = func(path string) ([]gentity.SubordinatedPage, error) {
					return []gentity.SubordinatedPage{
						{
							Path: "path/_test",
						},
					}, nil
				}
				problemRepo.EXPECT().GetAll().Return(nil, nil)

			},
		},
		{
			name: "問題が存在しない場合作成される",
			mockInit: func(growiClient *mockGrowiClient, problemRepo mock.MockProblemRepository) {
				growiClient.getSubordinatedPageReturn = func(path string) ([]gentity.SubordinatedPage, error) {
					return []gentity.SubordinatedPage{
						{
							Path: "path/test",
						},
					}, nil
				}
				growiClient.getPageReturn = func(path string) (*gentity.Page, error) {
					return &gentity.Page{
						Path: "path/test",
						Revision: gentity.Revision{
							Body: `---
code: ABC
title: test
point: 0
solvedCriterion: 0
type: normal
---
body`,
						},
					}, nil
				}
				problemRepo.EXPECT().GetAll().Return(nil, nil)
				problemRepo.EXPECT().Create(
					&entity.Problem{
						Code:     "ABC",
						AuthorID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Title:    "test",
						Type:     entity.NormalType,
						Body: `---
code: ABC
title: test
point: 0
solvedCriterion: 0
type: normal
---
body`,
					},
				).Return(nil, nil)
			},
		},
		{
			name: "問題が存在するかつ更新がない場合はスキップされる",
			mockInit: func(growiClient *mockGrowiClient, problemRepo mock.MockProblemRepository) {
				growiClient.getSubordinatedPageReturn = func(path string) ([]gentity.SubordinatedPage, error) {
					return []gentity.SubordinatedPage{
						{
							Path: "path/test",
						},
					}, nil
				}
				growiClient.getPageReturn = func(path string) (*gentity.Page, error) {
					return &gentity.Page{
						Path: "path/test",
						Revision: gentity.Revision{
							Body: `---
code: ABC
title: test
point: 0
solvedCriterion: 0
type: normal
---
body`,
						},
						UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				}
				problemRepo.EXPECT().GetAll().Return(
					[]*entity.Problem{
						{
							Base: entity.Base{
								ID:        uuid.MustParse("00000000-0000-4000-a000-000000000000"),
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
							Code:     "ABC",
							AuthorID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
							Title:    "test",
							Type:     entity.NormalType,
							Body:     `test`,
						},
					},
					nil)
				problemRepo.EXPECT().Create(gomock.Any()).Times(0)
				problemRepo.EXPECT().Update(gomock.Any()).Times(0)
			},
		},
		{
			name: "問題が存在するかつ更新がある場合は更新される",
			mockInit: func(growiClient *mockGrowiClient, problemRepo mock.MockProblemRepository) {
				growiClient.getSubordinatedPageReturn = func(path string) ([]gentity.SubordinatedPage, error) {
					return []gentity.SubordinatedPage{
						{
							Path: "path/test",
						},
					}, nil
				}
				growiClient.getPageReturn = func(path string) (*gentity.Page, error) {
					return &gentity.Page{
						Path: "path/test",
						Revision: gentity.Revision{
							Body: `---
code: ABC
title: test
point: 0
solvedCriterion: 0
type: normal
---
body`,
						},
						UpdatedAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
						CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				}
				problemRepo.EXPECT().GetAll().Return(
					[]*entity.Problem{
						{
							Base: entity.Base{
								ID:        uuid.MustParse("00000000-0000-4000-a000-000000000000"),
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
							Code:     "ABC",
							AuthorID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
							Title:    "test",
							Type:     entity.NormalType,
							Body:     `test`,
						},
					},
					nil)
				problemRepo.EXPECT().Update(
					&entity.Problem{
						Base: entity.Base{
							ID:        uuid.MustParse("00000000-0000-4000-a000-000000000000"),
							UpdatedAt: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
							CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Code:     "ABC",
						AuthorID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Title:    "test",
						Type:     entity.NormalType,
						Body: `---
code: ABC
title: test
point: 0
solvedCriterion: 0
type: normal
---
body`,
					},
				).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := mockGrowiClient{}
			problemRepo := mock.NewMockProblemRepository(ctrl)

			if tt.mockInit != nil {
				tt.mockInit(&client, *problemRepo)
			}

			s := NewGrowiProblemSyncService(
				client,
				"path",
				"00000000-0000-4000-a000-000000000000",
				problemRepo,
			)

			// when then
			if err := s.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockGrowiClient struct {
	getSubordinatedPageReturn func(path string) ([]gentity.SubordinatedPage, error)
	getPageReturn             func(path string) (*gentity.Page, error)
}

var _ growi_client.Client = (*mockGrowiClient)(nil)

func (c mockGrowiClient) Init() error {
	//TODO implement me
	panic("implement me")
}

func (c mockGrowiClient) GetSubordinatedPage(path string) ([]gentity.SubordinatedPage, error) {
	return c.getSubordinatedPageReturn(path)
}

func (c mockGrowiClient) GetPage(path string) (*gentity.Page, error) {
	return c.getPageReturn(path)
}
