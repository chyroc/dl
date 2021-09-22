package helper

func DownloadMp3(mp3 *MP3, chapter bool) (string, error) {
	file, err := Download(mp3.DownloadUrl, chapter)
	if err != nil {
		return "", err
	}

	return file, UpdateMp3Tag(mp3, &mp3.Tag, file)
}

func DownloadMp32(mp3 *MP3, chapter bool, target string) error {
	file, err := DownloadMp3(mp3, chapter)
	if err != nil {
		return err
	}
	return Rename(file, target)
}
