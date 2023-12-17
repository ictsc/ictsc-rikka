package entity

import (
	"database/sql/driver"
	"github.com/adrg/frontmatter"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

type Problem struct {
	Base

	Code              string             `json:"code" gorm:"unique;index"`
	Title             string             `json:"title"`
	Body              string             `json:"body,omitempty"`
	Type              ProblemType        `json:"type" yaml:"type" gorm:"type:enum('normal','multiple');default:'normal'"`
	CorrectAnswers    YAMLCorrectAnswers `json:"-" gorm:"type:text"`
	Point             uint               `json:"point"`
	PreviousProblemID *uuid.UUID         `json:"previous_problem_id"`
	PreviousProblem   *Problem           `json:"-"`
	SolvedCriterion   uint               `json:"solved_criterion"`
}

func (p *Problem) Validate() error {
	if matches, err := regexp.Match("[a-zA-Z]{3}", []byte(p.Code)); err != nil {
		return err
	} else if !matches {
		return errors.New("code must match the pattern [A-Z]{3}")
	}

	if !(p.SolvedCriterion <= p.Point) {
		return errors.New("solved_criterion must be less than or equal to point")
	}

	if !(p.Type == NormalType || p.Type == MultipleType) {
		return errors.New("invalid problem type")
	}

	if p.Type == MultipleType && len(p.CorrectAnswers) == 0 {
		return errors.New("multiple type problem must have at least one correct answer")
	}

	return nil
}

func (p *Problem) DeleteMatterQuestionWithQuestionFieldAttach() error {
	var matter = &ProblemFrontMatter{}
	body, err := frontmatter.Parse(strings.NewReader(p.Body), matter)
	if err != nil {
		return errors.Wrap(err, "failed to parse frontmatter")
	}
	if err := matter.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate frontmatter")
	}

	// matter から question を削除
	p.CorrectAnswers = matter.CorrectAnswers
	matter.CorrectAnswers = nil
	matterStr, err := matter.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode frontmatter")
	}
	p.Body = "---\n" + matterStr + "---\n" + string(body)

	return nil
}

type ProblemWithIsAnswered struct {
	Problem

	IsAnswered bool `json:"is_answered" gorm:"column:is_answered"`
}

type ProblemWithAnswerInformation struct {
	Problem

	Unchecked            uint `json:"unchecked"`
	UncheckedNearOverdue uint `json:"unchecked_near_overdue"`
	UncheckedOverdue     uint `json:"unchecked_overdue"`

	CurrentPoint uint `json:"current_point"`
	IsSolved     bool `json:"is_solved"`
}

type ProblemWithCurrentPoint struct {
	Problem

	IsAnswered   bool `json:"is_answered"`
	CurrentPoint uint `json:"current_point"`
	IsSolved     bool `json:"is_solved"`
}

type YAMLCorrectAnswers []CorrectAnswer

func (y *YAMLCorrectAnswers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal YAML")
	}

	return yaml.Unmarshal(bytes, y)
}

func (y YAMLCorrectAnswers) Value() (driver.Value, error) {
	if len(y) == 0 {
		return nil, nil
	}
	return yaml.Marshal(y)
}
