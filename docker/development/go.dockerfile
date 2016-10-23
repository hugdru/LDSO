FROM debian:stable
MAINTAINER "Hugo Drumond" hugdru@gmail.com

# https://hub.docker.com/_/golang/

ENV \
  USER="server" \
  GROUP="server" \
  GOLANG_VERSION="1.7.3" \
  GOLANG_DOWNLOAD_SHA256="508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e" \
  WATCHMAN_VERSION="v4.7.0" \
  PERSISTENT_APT_PACKAGES="git libpcre++0 libpcre3" \
  TEMPORARY_APT_PACKAGES="autoconf automake build-essential curl python-dev libpcre3-dev libpcre++-dev"

ENV HOME="/$USER"
ENV BACKEND_DIR="$HOME/backend"
ENV BACKEND_BINARY="$HOME/backend/bin/server"
ENV GOPATH="$HOME/go"
ENV HOST_ENTRYPOINT_DIR="docker/development/entrypoints"
ENV HOST_ENTRYPOINT_FILE="$HOST_ENTRYPOINT_DIR/go_entrypoint.sh"
ENV HOST_WATCHMAN_FILE="$HOST_ENTRYPOINT_DIR/go_watchman.sh"
ENV SCRIPTS_DIR="$HOME/bin"

RUN \
      apt-get update && apt-get upgrade -y && \
      apt-get install -y --no-install-recommends $PERSISTENT_APT_PACKAGES $TEMPORARY_APT_PACKAGES && \
      groupadd "$GROUP" && useradd -d "$HOME" -g "$GROUP" -s /bin/bash "$USER" && mkdir "$HOME" && \
      curl -fsSL "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz" -o golang.tar.gz && \
      echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - && \
      tar -C /usr/local -xzf golang.tar.gz && \
      rm golang.tar.gz && \
      git clone https://github.com/facebook/watchman.git /watchman && \
      cd watchman git checkout "$WATCHMAN_VERSION" && ./autogen.sh && ./configure && make && make install && \
      apt-get remove --purge -y $TEMPORARY_APT_PACKAGES && \
      rm -rf /var/lib/apt/lists/* /watchman && \
      apt-get autoremove -y && \
      apt-get clean all

VOLUME "$BACKEND_DIR"

WORKDIR "$SCRIPTS_DIR"
COPY "$HOST_ENTRYPOINT_FILE" entrypoint.sh
COPY "$HOST_WATCHMAN_FILE" watchman_command.sh

RUN chown -R "$USER":"$GROUP" "$HOME" && chmod +x "$SCRIPTS_DIR/entrypoint.sh" "$SCRIPTS_DIR/watchman_command.sh"

USER "$USER"
WORKDIR "$HOME"

ENV PATH="$GOPATH/bin:$SCRIPTS_DIR:/usr/local/go/bin:$PATH"
RUN go get github.com/constabulary/gb/...

EXPOSE 8080

WORKDIR "$BACKEND_DIR"
ENTRYPOINT ["entrypoint.sh"]
