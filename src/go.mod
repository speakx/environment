module environment

go 1.13

replace consolespeakx => /Users/Kim/Documents/code/github.kimj/research/consolespeakx/src

replace go-simple => /Users/Kim/Documents/code/github.kimj/research/go-simple/src

replace grpc-flatc => /Users/Kim/Documents/code/github.kimj/research/grpc-flatc/src

replace grpc-idl => /Users/Kim/Documents/code/github.kimj/research/grpc-idl/src

replace rocksdb => /Users/Kim/Documents/code/github.kimj/research/rocksdb/src

replace src/environment => ../../pkg/src/environment

require (
	github.com/edsrzf/mmap-go v1.0.0
	golang.org/x/sys v0.0.0-20191110163157-d32e6e3b99c4 // indirect
)
