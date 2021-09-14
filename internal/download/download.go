package download

type Downloader interface {
	Title() string
	Download() error
}
