FROM debian:stable
MAINTAINER "Hugo Drumond" hugdru@gmail.com

# https://hub.docker.com/_/golang/

ENV \
  USER="server" \
  GROUP="server" \
  GOLANG_VERSION="1.7.3" \
  GOLANG_DOWNLOAD_SHA256="508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e" \
  PERSISTENT_APT_PACKAGES="git ca-certificates" \
  TEMPORARY_APT_PACKAGES="curl" \
  DOCKERIZE_VERSION="v0.2.0" \
  DOCKERIZE_DOWNLOAD_SHA256="c0e2e33cfe066036941bf8f2598090bd8e01fdc05128490238b2a64cf988ecfb"

ENV HOME="/$USER"

ENV \
  SCRIPTS_DIR="$HOME/bin" \
  BACKEND_DIR="$HOME/backend" \
  HOST_ENTRYPOINT_DIR="docker/production/entrypoints"

ENV HOST_ENTRYPOINT_FILE="$HOST_ENTRYPOINT_DIR/go_entrypoint.sh"

RUN \
      apt-get update && apt-get upgrade -y && apt-get install -y --no-install-recommends \
      $PERSISTENT_APT_PACKAGES \
      $TEMPORARY_APT_PACKAGES && \
      groupadd "$GROUP" && useradd -d "$HOME" -g "$GROUP" -s /bin/bash "$USER" && mkdir "$HOME" && \
      curl -fsSL "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz" -o golang.tar.gz && \
      echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - && \
      tar -C /usr/local -xzf golang.tar.gz && \
      rm golang.tar.gz && \
      curl -fsSL "https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz" -o dockerize.tar.gz && \
      echo "$DOCKERIZE_DOWNLOAD_SHA256  dockerize.tar.gz" | sha256sum -c - && \
      tar -C /usr/local/bin -xzvf dockerize.tar.gz && rm dockerize.tar.gz && \
      apt-get remove -y $TEMPORARY_APT_PACKAGES && \
      rm -rf /var/lib/apt/lists/* && \
      apt-get autoremove -y && \
      apt-get clean all

COPY "$HOST_ENTRYPOINT_FILE" "$SCRIPTS_DIR/entrypoint.sh"
COPY backend "$BACKEND_DIR"
RUN chown -R "$USER":"$GROUP" "$HOME" && chmod +x "$SCRIPTS_DIR/entrypoint.sh"

USER "$SERVER"
WORKDIR "$HOME"

ENV GOPATH="$HOME/go"
ENV PATH="$GOPATH/bin:$SCRIPTS_DIR:/usr/local/go/bin:$PATH"
RUN \
      go get github.com/constabulary/gb/... && \
      git config --global http.https://gopkg.in.followRedirects true && \
      cd backend && gb vendor restore && gb build server && \
      find . -maxdepth 1 ! \( -name 'bin' -o -name "." -o -name ".." \) -exec rm -rf {} +

EXPOSE 8080

WORKDIR "$BACKEND_DIR"
ENTRYPOINT ["entrypoint.sh"]
