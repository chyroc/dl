package helper

import (
	"net/http"

	"github.com/bogem/id3v2"
	"github.com/chyroc/dl/internal/config"
)

type Tag struct {
	Title      string
	Artist     string
	Album      string
	Year       string
	Track      string
	CoverImage string
}

type MP3 struct {
	FileName    string
	SavePath    string
	Playable    bool
	DownloadUrl string
	Tag         Tag
	Origin      int
}

func UpdateMp3Tag(mp3 *MP3, tag *Tag, file string) error {
	data, err := config.ReqCli.New(http.MethodGet, tag.CoverImage).Bytes()
	if err != nil {
		return err
	}

	tagID, err := id3v2.Open(file, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tagID.Close()

	tagID.SetDefaultEncoding(id3v2.EncodingUTF8)
	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpg",
		PictureType: id3v2.PTOther,
		Picture:     data,
	}
	tagID.AddAttachedPicture(pic)
	tagID.SetTitle(mp3.Tag.Title)
	tagID.SetArtist(mp3.Tag.Artist)
	tagID.SetAlbum(mp3.Tag.Album)
	tagID.SetYear(mp3.Tag.Year)
	textFrame := id3v2.TextFrame{
		Encoding: id3v2.EncodingUTF8,
		Text:     mp3.Tag.Track,
	}
	tagID.AddFrame(tagID.CommonID("Track number/Position in set"), textFrame)

	return tagID.Save()
}
