package ports

type DocumentParser interface {
	DownloadPage(location string) error
	GetDocumentVersion() string
	GetTitle() string
	GetExternalLinkCount() int
	GetInternalLinkCount() int
	GetContainsLogin() bool
	GetHeaderOneCount() int
	GetHeaderTwoCount() int
	GetHeaderThreeCount() int
	GetHeaderFourCount() int
	GetHeaderFiveCount() int
	GetHeaderSixCount() int
}
