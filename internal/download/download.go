package download

type Downloader interface {
	Title() string
	TargetFile() string
	Download() error
}
