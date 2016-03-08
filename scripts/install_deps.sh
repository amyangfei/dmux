#!/bin/bash

install_remote_dep() {
    go get -u -v github.com/astaxie/beego
    go get -u -v github.com/koding/multiconfig
}

install_local_dep() {
    cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
    bash +x $cur/dev_deps_update.sh
}

install_remote_dep $*
install_local_dep $*
