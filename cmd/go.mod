module photo-app

go 1.18

require (
	github.com/dapr/go-sdk v1.4.0
	models v0.1.0 // indirect
	utils v0.1.0
)

require (
	github.com/dapr/dapr v1.6.1-0.20220223052122-c503710ab2d1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.0.0-20220601225756-64ec528b34cd // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210831024726-fe130286e0e2 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace models v0.1.0 => ./models

replace utils v0.1.0 => ./utils
