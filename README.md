# dl - download.

## Feature

download website video.

support:

- v.youku.com
- haokan.baodu.com
- video.sina.com.cn

## Install

### By Go Get

```shell
go get github.com/chyroc/dl
```

### By Brew

```shell
brew tap chyroc/tap
brew install chyroc/tap/dl
```

## Usage

```shell
dl 'https://haokan.baidu.com/v?vid=7249594116085322255'

[meta] haokan.baidu.com
[download] 蜡笔小新：小新一家来采茶园玩，猜茶环节小新胜出，主办方亏大了,动漫,日本动漫,好看视频.mp4
[download] 13.65 MiB / 13.65 MiB 100 % [==================] 298.11 KiB/s
```

# Other

My wife is a teacher and needs to play some videos during the teaching process,
so I made a tool to help her download these videos.

Why not youtube-dl, because these video sites may be unique to China.