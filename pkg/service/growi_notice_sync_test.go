package service

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	gentity "github.com/ictsc/growi_client/entity"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mock"
	"testing"
	"time"
)

func TestNewGrowiNoticeSyncService(t *testing.T) {
	tests := []struct {
		name     string
		mockInit func(growiClient *mockGrowiClient, noticeRepo mock.MockNoticeRepository)
		wantErr  bool
	}{
		{
			name: "ページ名が _ で始まるページがスキップされる",
			mockInit: func(growiClient *mockGrowiClient, noticeRepo mock.MockNoticeRepository) {
				growiClient.getSubordinatedPageReturn = func(path string) ([]gentity.SubordinatedPage, error) {
					return []gentity.SubordinatedPage{
						{
							Path: "path/_test",
						},
					}, nil
				}
				noticeRepo.EXPECT().GetAll().Return(nil, nil)
			},
		},
		{
			name: "お知らせが存在しない場合作成される",
			mockInit: func(growiClient *mockGrowiClient, noticeRepo mock.MockNoticeRepository) {
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
title: test
---
body`,
						},
					}, nil
				}
				noticeRepo.EXPECT().GetAll().Return(nil, nil)
				noticeRepo.EXPECT().Create(
					&entity.Notice{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						},
						SourceId: "test",
						Title:    "test",
						Body:     "body",
					},
				).Return(nil, nil)
			},
		},
		{
			name: "お知らせが存在する場合更新される",
			mockInit: func(growiClient *mockGrowiClient, noticeRepo mock.MockNoticeRepository) {
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
title: test
---
body`,
						},
						UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				}
				noticeRepo.EXPECT().GetAll().Return([]*entity.Notice{
					{
						Base: entity.Base{
							ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
							UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						SourceId: "test",
						Title:    "test",
						Body:     "body",
					},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := mockGrowiClient{}
			noticeRepo := mock.NewMockNoticeRepository(ctrl)

			if tt.mockInit != nil {
				tt.mockInit(&client, *noticeRepo)
			}

			s := NewGrowiNoticeSyncService(
				client,
				"path",
				noticeRepo,
			)

			// when then
			if err := s.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
