#!/bin/bash

org=${PWD%/*}
org=${org##*/}
repository=${PWD##*/}
echo "** org:$org"
echo "** repository:$repository"
echo

echo "rebuild proto"
sh ./gen-proto.sh
echo "cleanup old mod folder:$repository -> ../pkg/src/$org/repository"
rm -rf ../pkg/src/$org
echo "copy mod files:$repository -> ../pkg/src/$org/$repository"
echo "mkdir -p ../pkg/src/$org/$repository"
mkdir -p ../pkg/src/$org/$repository
echo "cp -r ./src/* ../pkg/src/$org/$repository"
cp -r ./src/* ../pkg/src/$org/$repository
echo

dir=$(pwd)
echo "create go mod:$repository"
cd ../pkg/src/$org/$repository
rm -f go.mod
rm -f go.sum
go mod init $repository
cd $dir