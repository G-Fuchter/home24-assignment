package ports

type DocumentParser interface {
	DownloadDocument(location string) error
	GetDocumentVersion() (string, error)
	GetTitle() (string, error)
	GetExternalLinkCount() (int, error)
	GetInternalLinkCount() (int, error)
	GetContainsLogin() (bool, error)
	GetHeaderOneCount() (int, error)
	GetHeaderTwoCount() (int, error)
	GetHeaderThreeCount() (int, error)
	GetHeaderFourCount() (int, error)
	GetHeaderFiveCount() (int, error)
	GetHeaderSixCount() (int, error)
}
