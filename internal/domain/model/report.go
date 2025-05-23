package model

type WebPageReport struct {
	DocumentVersion   string
	Title             string
	ExternalLinkCount int
	InternalLinkCount int
	ContainsLogin     bool
	HeaderOneCount    int
	HeaderTwoCount    int
	HeaderThreeCount  int
	HeaderFourCount   int
	HeaderFiveCount   int
	HeaderSixCount    int
}
