#!/bin/bash

# install dependcies for local developing

cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
src_dst=$GOPATH/src/github.com/amyangfei/dmux

mkdir -p $src_dst
rm -rf $src_dst/*

modules=('utils' 'store' 'registry/models' 'registry/controllers')

for module in ${modules[@]}; do
    mkdir -p $src_dst/$module
    cp -r $cur/../$module $src_dst/$module/..
    cd $src_dst/$module
    echo "installing dmux/$module"
    go install
done

echo "done!"
