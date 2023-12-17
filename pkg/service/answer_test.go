package service

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository/mock"
	"github.com/pkg/errors"
	"testing"
)

type answerMatcher struct {
	expected *entity.Answer
}

// Matches *entity.Answer のマッチャー
func (a answerMatcher) Matches(x interface{}) bool {
	arg, ok := x.(*entity.Answer)
	if !ok {
		return false
	}
	if a.expected.UserGroupID != arg.UserGroupID {
		return false
	}
	if a.expected.ProblemID != arg.ProblemID {
		return false
	}
	if *a.expected.Point != *arg.Point {
		return false
	}
	if a.expected.Body != arg.Body {
		return false
	}
	return true
}

func (a answerMatcher) String() string {
	return fmt.Sprintf("is equal to %v", a.expected)
}

func TestAnswerService_Create(t *testing.T) {
	type args struct {
		req *CreateAnswerRequest
	}
	tests := []struct {
		name     string
		args     args
		mockInit func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository)
		want     *entity.Answer
		wantErr  bool
	}{
		{
			name: "UserGroup の取得に失敗する",
			args: args{
				req: &CreateAnswerRequest{
					UserGroup: &entity.UserGroup{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						},
					},
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository) {
				answerRepo.EXPECT().FindByUserGroup(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "選択式問題（ラジオボタン）が採点され、Answer が作成される",
			args: args{
				req: &CreateAnswerRequest{
					UserGroup: &entity.UserGroup{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						},
					},
					Body:      `[{"group": 0, "value": [3], "type": "radio"}]`,
					ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository) {
				answerRepo.EXPECT().FindByUserGroup(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(nil, nil).Times(1)
				problemRepo.EXPECT().FindByID(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(
					&entity.Problem{
						Code: "ABC",
						CorrectAnswers: []entity.CorrectAnswer{
							{
								Type:   entity.CheckBox,
								Column: []uint{3},
								Scoring: entity.Scoring{
									Correct: 10,
									PartialCorrect: func() *uint {
										val := uint(5)
										return &val
									}(),
								},
							},
						},
						Type: entity.MultipleType,
					},
					nil,
				).Times(1)
				answerRepo.EXPECT().Create(answerMatcher{
					expected: &entity.Answer{
						UserGroupID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Point: func() *uint {
							val := uint(10)
							return &val
						}(),
						Body:      `[{"group": 0, "value": [3], "type": "radio"}]`,
						ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
					},
				},
				).Return(&entity.Answer{}, nil).Times(1)
			},
			want: &entity.Answer{},
		},
		{
			name: "選択式問題（チェックボックス）かつ一文間違った状態で採点され、Answer が作成される",
			args: args{
				req: &CreateAnswerRequest{
					UserGroup: &entity.UserGroup{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						},
					},
					Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"}]`,
					ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository) {
				answerRepo.EXPECT().FindByUserGroup(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(nil, nil).Times(1)
				problemRepo.EXPECT().FindByID(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(
					&entity.Problem{
						Code: "ABC",
						CorrectAnswers: []entity.CorrectAnswer{
							{
								Type:   entity.CheckBox,
								Column: []uint{0, 1, 2},
								Scoring: entity.Scoring{
									Correct: 10,
									PartialCorrect: func() *uint {
										val := uint(2)
										return &val
									}(),
								},
							},
						},
						Type: entity.MultipleType,
					},
					nil,
				).Times(1)
				answerRepo.EXPECT().Create(answerMatcher{
					expected: &entity.Answer{
						UserGroupID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Point: func() *uint {
							val := uint(4)
							return &val
						}(),
						Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"}]`,
						ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
					},
				},
				).Return(&entity.Answer{}, nil).Times(1)
			},
			want: &entity.Answer{},
		},
		{
			name: "選択式問題（チェックボックス）が採点され、Answer が作成される",
			args: args{
				req: &CreateAnswerRequest{
					UserGroup: &entity.UserGroup{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						},
					},
					Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"}]`,
					ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository) {
				answerRepo.EXPECT().FindByUserGroup(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(nil, nil).Times(1)
				problemRepo.EXPECT().FindByID(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(
					&entity.Problem{
						Code: "ABC",
						CorrectAnswers: []entity.CorrectAnswer{
							{
								Type:   entity.CheckBox,
								Column: []uint{0, 1},
								Scoring: entity.Scoring{
									Correct: 10,
									PartialCorrect: func() *uint {
										val := uint(5)
										return &val
									}(),
								},
							},
						},
						Type: entity.MultipleType,
					},
					nil,
				).Times(1)
				answerRepo.EXPECT().Create(answerMatcher{
					expected: &entity.Answer{
						UserGroupID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Point: func() *uint {
							val := uint(10)
							return &val
						}(),
						Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"}]`,
						ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
					},
				},
				).Return(&entity.Answer{}, nil).Times(1)
			},
			want: &entity.Answer{},
		},
		{
			name: "ラジオボタン問題とチェックボックス問題が混ざった状態で採点され、Answer が作成される",
			args: args{
				req: &CreateAnswerRequest{
					UserGroup: &entity.UserGroup{
						Base: entity.Base{
							ID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						},
					},
					Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"},{"group": 2, "value": [3], "type": "radio"}]`,
					ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				},
			},
			mockInit: func(userRepo mock.MockUserRepository, answerRepo mock.MockAnswerRepository, problemRepo mock.MockProblemRepository) {
				answerRepo.EXPECT().FindByUserGroup(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(nil, nil).Times(1)
				problemRepo.EXPECT().FindByID(
					uuid.MustParse("00000000-0000-4000-a000-000000000000"),
				).Return(
					&entity.Problem{
						Code: "ABC",
						CorrectAnswers: []entity.CorrectAnswer{
							{
								Type:   entity.CheckBox,
								Column: []uint{0, 1, 2},
								Scoring: entity.Scoring{
									Correct: 10,
									PartialCorrect: func() *uint {
										val := uint(2)
										return &val
									}(),
								},
							},
							{
								Type:   entity.CheckBox,
								Column: []uint{0, 1, 2},
								Scoring: entity.Scoring{
									Correct: 100,
									PartialCorrect: func() *uint {
										val := uint(50)
										return &val
									}(),
								},
							},
							{
								Type:   entity.RadioButton,
								Column: []uint{3},
								Scoring: entity.Scoring{
									Correct: 10,
								},
							},
						},
						Type: entity.MultipleType,
					},
					nil,
				).Times(1)
				answerRepo.EXPECT().Create(answerMatcher{
					expected: &entity.Answer{
						UserGroupID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
						Point: func() *uint {
							val := uint(14)
							return &val
						}(),
						Body:      `[{"group": 0, "value": [0, 1], "type": "checkbox"},{"group": 2, "value": [3], "type": "radio"}]`,
						ProblemID: uuid.MustParse("00000000-0000-4000-a000-000000000000"),
					},
				},
				).Return(&entity.Answer{}, nil).Times(1)
			},
			want: &entity.Answer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			answerRepo := mock.NewMockAnswerRepository(ctrl)
			problemRepo := mock.NewMockProblemRepository(ctrl)

			s := NewAnswerService(
				false,
				0,
				"",
				userRepo,
				answerRepo,
				problemRepo,
			)

			if tt.mockInit != nil {
				tt.mockInit(*userRepo, *answerRepo, *problemRepo)
			}

			// when
			got, err := s.Create(tt.args.req)

			// then
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
