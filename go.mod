module github.com/JungleMC/java-edition

go 1.16

replace github.com/JungleMC/sdk => ../sdk

require (
	github.com/JungleMC/sdk v0.0.0-20210810042112-e30cdbe2f38a
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-redis/redis/v8 v8.11.2
	github.com/google/uuid v1.1.2
	github.com/stretchr/testify v1.7.0 // indirect
	google.golang.org/protobuf v1.27.1
)
