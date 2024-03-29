all: generate build

generate:
	( cd .github/cmd/generate_register && go build -o /tmp/dl_generate_register main.go && cd ../../.. && /tmp/dl_generate_register )
	( cd .github/cmd/generate_readme && go build -o /tmp/dl_generate_readme main.go && cd ../../.. && /tmp/dl_generate_readme )

build:
	CGO_ENABLED=1 go build -o /tmp/dl main.go