#!/usr/bin/env bash

set -e

cd /

# node
curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
sudo apt-get install -y nodejs

# go
curl -fsSL "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz" -o golang.tar.gz
echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c -
tar -C /usr/local -xzf golang.tar.gz
rm golang.tar.gz
export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"

