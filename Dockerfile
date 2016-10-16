FROM debian:stable
MAINTAINER "Hugo Drumond" hugdru@gmail.com

# https://hub.docker.com/_/golang/
# https://hub.docker.com/_/node/

ENV \
  USER="server" \
  GROUP="server" \
  GOLANG_VERSION="1.7.1" \
  GOLANG_DOWNLOAD_SHA256="43ad621c9b014cde8db17393dc108378d37bc853aa351a6c74bf6432c1bbd182" \
  NPM_CONFIG_LOGLEVEL="silent" \
  NODE_VERSION="6.8.1"

ENV HOME="/$USER"

RUN \
      apt-get update && apt-get upgrade -y && apt-get install -y \
      curl \
      git \
      xz-utils && \
      groupadd "$GROUP" && useradd -d "$HOME" -g "$GROUP" -s /bin/bash "$USER" && mkdir "$HOME" && \
      curl -fsSL "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz" -o golang.tar.gz && \
      echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - && \
      tar -C /usr/local -xzf golang.tar.gz && \
      rm golang.tar.gz && \
      set -ex && \
      for key in \
        9554F04D7259F04124DE6B476D5A82AC7E37093B \
        94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
        0034A06D9D9B0064CE8ADF6BF1747F4AD2306D93 \
        FD3A5288F042B6850C66B31F09FE44734EB7990E \
        71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
        DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
        C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
        B9AE9905FFD7803F25714661B63B535A4C206CA9 \
      ; do \
        gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$key"; \
      done && \
      curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz" && \
      curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" && \
      gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc && \
      grep "node-v$NODE_VERSION-linux-x64.tar.xz\$" SHASUMS256.txt | sha256sum -c - && \
      mkdir /usr/local/node && \
      tar -xJf "node-v$NODE_VERSION-linux-x64.tar.xz" -C /usr/local/node --strip-components=1 && \
      rm "node-v$NODE_VERSION-linux-x64.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt && \
      apt-get remove curl xz-utils -y && \
      rm -rf /var/lib/apt/lists/* && \
      apt-get autoremove -y && \
      apt-get clean all

COPY backend /server/backend
COPY frontend /server/frontend
RUN chown -R "$USER":"$GROUP" "$HOME"

USER server
WORKDIR /server

ENV GOPATH="$HOME/go"
ENV PATH="$GOPATH/bin:/usr/local/go/bin:/usr/local/node/bin:$PATH"
ENV PREFIX="$HOME/packages"
ENV PATH="$PREFIX/bin:$PATH"
RUN \
      go get github.com/constabulary/gb/... && \
      mkdir -p "$PREFIX/bin" && npm install -g npm && npm install -g yarn && \
      yarn global add angular-cli && \
      cd frontend && npm install && ng build -prod && \
      find . -maxdepth 1 ! \( -name 'dist' -o -name '.' -o -name '..' \) -exec rm -rf {} + && \
      cd ../backend && gb vendor restore && gb build server && \
      rm -rf pkg src vendor && \
      cd ~ && rm -rf packages .npm .yarn .yarn-cache .gnupg .config go

EXPOSE 80

WORKDIR /server/backend/
ENTRYPOINT ["./bin/server"]
