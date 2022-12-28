package helper

import (
	"os"
	"os/exec"
)

func M3u8ToMp4(file string) (string, error) {
	ffCli, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", err
	}

	// ffmpeg -i ./iii -vcodec copy -acodec copy -map 0:v -map 0:a ./2222.mp4
	f, err := os.CreateTemp("", "dl-m3u8-output-*.mp4")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(ffCli, "-y", "-i", file, "-vcodec", "copy", "-acodec", "copy", "-map", "0:v", "-map", "0:a", f.Name())
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
