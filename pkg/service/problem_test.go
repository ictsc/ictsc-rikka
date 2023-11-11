package service

import (
	"github.com/golang/mock/gomock"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mock"
	"reflect"
	"testing"
)

func TestProblemService_Create(t *testing.T) {
	type args struct {
		req *CreateProblemRequest
	}
	tests := []struct {
		name     string
		args     args
		mockInit func(userRepo mock.MockUserRepository, problemRepo mock.MockProblemRepository, answerRepo mock.MockAnswerRepository)
		want     *entity.Problem
		wantErr  bool
	}{
		{
			name: "問題のバリデーションに引っかかりエラーになる",
			args: args{
				req: &CreateProblemRequest{
					Code: "code",
				},
			},
			wantErr: true,
		},
		{
			name: "matter が正しくパースできずエラーになる",
			args: args{
				req: &CreateProblemRequest{
					Code: "ABC",
					Body: "invalid matter",
				},
			},
			wantErr: true,
		},
		{
			name: "問題の作成に成功する",
			args: args{
				req: &CreateProblemRequest{
					Code: "ABC",
					Body: `---
code: ABC
title: title
point: 100
solvedCriterion: 50
type: normal
---
body`,
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, problemRepo mock.MockProblemRepository, answerRepo mock.MockAnswerRepository) {
				problemRepo.EXPECT().Create(
					&entity.Problem{
						Code: "ABC",
						Body: `---
code: ABC
title: title
point: 100
solvedCriterion: 50
type: normal
---
body`,
					},
				).Return(&entity.Problem{}, nil).Times(1)
			},
			want: &entity.Problem{},
		},
		{
			name: "問題の作成に成功するかつ question を削除する",
			args: args{
				req: &CreateProblemRequest{
					Code: "ABC",
					Body: `---
code: ABC
title: title
point: 100
solvedCriterion: 50
type: multiple
questions:
  - type: radio
    correct_answers:
      - 1
    scoring:
      correct: 1
---
body`,
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, problemRepo mock.MockProblemRepository, answerRepo mock.MockAnswerRepository) {
				var correct uint = 1
				problemRepo.EXPECT().Create(
					&entity.Problem{
						Code: "ABC",
						Body: `---
code: ABC
title: title
point: 100
solvedCriterion: 50
type: multiple
---
body`,
						Answer: []entity.Question{
							{
								Type:           entity.RadioButton,
								CorrectAnswers: []uint{1},
								Scoring: entity.Scoring{
									Correct: correct,
								},
							},
						},
					},
				).Return(&entity.Problem{}, nil).Times(1)
			},
			want: &entity.Problem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			problemRepo := mock.NewMockProblemRepository(ctrl)
			answerRepo := mock.NewMockAnswerRepository(ctrl)

			s := NewProblemService(
				0,
				userRepo,
				problemRepo,
				answerRepo,
			)

			if tt.mockInit != nil {
				tt.mockInit(*userRepo, *problemRepo, *answerRepo)
			}

			// when
			got, err := s.Create(tt.args.req)

			// then
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
