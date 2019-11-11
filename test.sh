#!/bin/bash

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# 测试选项
if [ ! -n "$1" ] ;then
    echo "you need input test target { all | byteio | log | mmap | rocksdb }."
    exit
else
    echo "the test target is $1"
    echo
fi
target=$1
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

org=${PWD%/*}
org=${org##*/}
repository=${PWD##*/}
echo "** org:$org"
echo "** repository:$repository"
echo

# 重新造一遍 go mod
sh ./shell/gen-proto.sh
# sh ./shell/configure.sh
sh ./publish.sh # 把代码发布到pkg目录，如果test的文件对其他目录有依赖则通过gopath引入
curdir=$(pwd)
pkgenv=${curdir%/*}"/pkg"
export GOPATH=$GOPATH:${pkgenv}
echo $GOPATH

# byteio test
if [ "$target" == "all" ] || [ "$target" == "byteio" ] ;then
    go test -v -bench=".*" ./src/byteio/byteio_test.go ./src/byteio/byteio.go
    go test -bench=".*" ./src/byteio/byteio_test.go ./src/byteio/byteio.go
fi

# log test
if [ "$target" == "all" ] || [ "$target" == "log" ] ;then
    go test -v -bench=".*" ./src/logger/log_test.go ./src/logger/log.go
    rm -f ./src/logger/test.log*
fi

# mmap cache
if [ "$target" == "all" ] || [ "$target" == "mmap" ] ;then
    # go get github.com/edsrzf/mmap-go
    go test -v ./src/mmapcache/mmapcache_test.go ./src/mmapcache/mmapcachepool.go ./src/mmapcache/mmapcache.go
    go test -bench=".*" ./src/mmapcache/mmapcache_test.go ./src/mmapcache/mmapcachepool.go ./src/mmapcache/mmapcache.go
fi

# rocksdb test
if [ "$target" == "all" ] || [ "$target" == "rocksdb" ] ;then
    go get github.com/tecbot/gorocksdb
    go test -v ./src/rocksdbimp/rocksdbimp_test.go ./src/rocksdbimp/rocksdbimp.go
    go test -bench=".*" ./src/rocksdbimp/rocksdbimp_test.go ./src/rocksdbimp/rocksdbimp.go
fi
