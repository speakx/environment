#!/bin/bash

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# 测试选项
if [ ! -n "$1" ] ;then
    echo "you need input test target { all | log | srvinstance | transid | rocksdb }."
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
sh ./shell/configure.sh

# log test
if [ "$target" == "all" ] || [ "$target" == "log" ] ;then
    cd ./src
    go test -v -bench=".*" ./logger/log_test.go ./logger/log.go
    rm -f ./logger/test.log*
fi

# srvinstance test
if [ "$target" == "all" ] || [ "$target" == "srvinstance" ] ;then
    cd ./src
    go test -v ./srvinstance/tool_test.go ./srvinstance/tool.go
fi

# transid test
if [ "$target" == "all" ] || [ "$target" == "transid" ] ;then
    cd ./src
    go test -v -bench=".*" ./transid/transid_test.go
fi

# rocksdb test
if [ "$target" == "all" ] || [ "$target" == "rocksdb" ] ;then
    go get github.com/tecbot/gorocksdb
    cd ./src
    go test -v ./rocksdbimp/rocksdbimp_test.go ./rocksdbimp/rocksdbimp.go
    go test -v -bench=".*" ./rocksdbimp/rocksdbimp_test.go ./rocksdbimp/rocksdbimp.go
fi
