module environment

go 1.13

replace mmapcache => ../../mmapcache/src

replace single => ../../single/src

replace svrdemo => ../../svrdemo/src

require (
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tecbot/gorocksdb v0.0.0-20191019123150-400c56251341
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.5
)
