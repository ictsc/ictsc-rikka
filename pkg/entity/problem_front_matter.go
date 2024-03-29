package entity

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

type ProblemType string

const (
	NormalType   ProblemType = "normal"
	MultipleType ProblemType = "multiple"
)

type QuestionType string

const (
	RadioButton QuestionType = "radio"
	CheckBox    QuestionType = "checkbox"
)

type ProblemFrontMatter struct {
	Code              string          `yaml:"code"`
	Title             string          `yaml:"title"`
	Point             uint            `yaml:"point"`
	PreviousProblemID *string         `yaml:"previousProblemId,omitempty"`
	SolvedCriterion   uint            `yaml:"solvedCriterion"`
	Type              ProblemType     `yaml:"type"`
	CorrectAnswers    []CorrectAnswer `yaml:"correct_answers,omitempty"`
	ConnectInfo       []ConnectInfo   `yaml:"connectInfo,omitempty"`
}

type CorrectAnswer struct {
	Type    QuestionType `yaml:"type"`
	Column  []uint       `yaml:"column"`
	Scoring Scoring      `yaml:"scoring"`
}

type Scoring struct {
	Correct        uint  `yaml:"correct,omitempty"`
	PartialCorrect *uint `yaml:"partial_correct,omitempty"`
}

type ConnectInfo struct {
	Hostname string `yaml:"hostname"`
	Command  string `yaml:"command"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Type     string `yaml:"type"`
}

func (p *ProblemFrontMatter) Validate() error {
	if matches, err := regexp.Match("[a-zA-Z]{3}", []byte(p.Code)); err != nil {
		return err
	} else if !matches {
		return errors.New("code must match the pattern [A-Z]{3}")
	}

	if !(p.SolvedCriterion <= p.Point) {
		return errors.New("solved_criterion must be less than or equal to point")
	}

	if p.Type != NormalType && p.Type != MultipleType {
		return errors.New("invalid problem type")
	}

	if p.Type == MultipleType && len(p.CorrectAnswers) < 1 {
		return errors.New("multiple type problem must have at least one correct answer")
	}

	for _, ca := range p.CorrectAnswers {
		if ca.Type != RadioButton && ca.Type != CheckBox {
			return errors.New("invalid ca type")
		}

		if err := validateScoring(ca.Scoring); err != nil {
			return err
		}
	}

	return nil
}

func (p *ProblemFrontMatter) Encode() (string, error) {
	var buf strings.Builder

	enc := yaml.NewEncoder(&buf)

	if err := enc.Encode(p); err != nil {
		return "", err
	}
	defer enc.Close()

	return buf.String(), nil
}

func validateScoring(s Scoring) error {
	if s.PartialCorrect != nil && s.Correct < *s.PartialCorrect {
		return errors.New("partial_correct must be less than or equal to correct")
	}

	return nil
}
