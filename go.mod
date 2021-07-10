module github.com/junglemc/Service-JavaEditionHost

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/junglemc/Service-PlayerProvider v0.0.0-20210710163249-5059e486f2d7 // indirect
	github.com/junglemc/Service-StatusProvider v0.0.0-20210710134346-2f444d18fe62
	github.com/junglemc/nbt v1.2.0
	google.golang.org/grpc v1.39.0
)

replace github.com/junglemc/Service-PlayerProvider => ../Service-PlayerProvider
replace github.com/junglemc/Service-StatusProvider => ../Service-StatusProvider
