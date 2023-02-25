package entity

type FrontMatter struct {
	Code              string  `yaml:"code"`
	Title             string  `yaml:"title"`
	Point             uint    `yaml:"point"`
	PreviousProblemID *string `yaml:"previousProblemId"`
	SolvedCriterion   uint    `yaml:"solvedCriterion"`
	AuthorId          string  `yaml:"authorId"`
}
