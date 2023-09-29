# cmd/agent

To run with external config JSON file please use the following:
```
go run main.go -config="./configs/agent_config.json" 
```

To run with link flags please use the following commands:

*for zsh*
```
go run -ldflags "-X main.buildVersion=v19.0 -X main.buildDate=$(date +%d.%m.%Y) -X main.buildCommit=$(git rev-parse HEAD)" main.go
```

*for bash*
```
go run -ldflags "-X main.buildVersion=v19.0 -X 'main.buildDate=$(date +'%d.%m.%Y')' -X 'main.buildCommit=$(git rev-parse HEAD)'" main.go
```

To build with link flags please use the following commands:

*for zsh*
```
go build -ldflags "-X main.buildVersion=v19.0 -X main.buildDate=$(date +%d.%m.%Y) -X main.buildCommit=$(git rev-parse HEAD)"
```

*for bash*
```
go build -ldflags "-X main.buildVersion=v19.0 -X 'main.buildDate=$(date +'%d.%m.%Y')' -X 'main.buildCommit=$(git rev-parse HEAD)'"
```