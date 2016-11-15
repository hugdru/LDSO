FROM nginx:stable
MAINTAINER "Hugo Drumond" hugdru@gmail.com

# https://hub.docker.com/_/node/

ARG BUILD
ENV BUILD ${BUILD:-prod}

ENV \
  USER="server" \
  GROUP="server" \
  NPM_CONFIG_LOGLEVEL="info" \
  NODE_VERSION="6.9.1" \
  PERSISTENT_APT_PACKAGES="git ca-certificates python make g++" \
  TEMPORARY_APT_PACKAGES="curl xz-utils" \
  DOCKERIZE_VERSION="v0.2.0" \
  DOCKERIZE_DOWNLOAD_SHA256="c0e2e33cfe066036941bf8f2598090bd8e01fdc05128490238b2a64cf988ecfb" \
  YARN_VERSION="0.16.1"

ENV HOME="/$USER"
ENV FRONTEND_DIR="$HOME/frontend"
ENV HTDOCS="$FRONTEND_DIR/dist"
ENV PREFIX="$HOME/packages"

RUN \
      apt-get update && apt-get upgrade -y && \
      apt-get install -y --no-install-recommends $PERSISTENT_APT_PACKAGES $TEMPORARY_APT_PACKAGES && \
      groupadd "$GROUP" && useradd -d "$HOME" -g "$GROUP" -s /bin/bash "$USER" && mkdir "$HOME" && \
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
      curl -fsSL "https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz" -o dockerize.tar.gz && \
      echo "$DOCKERIZE_DOWNLOAD_SHA256  dockerize.tar.gz" | sha256sum -c - && \
      tar -C /usr/local/bin -xzvf dockerize.tar.gz && rm dockerize.tar.gz && \
      apt-get remove --purge -y $TEMPORARY_APT_PACKAGES && \
      rm -rf /var/lib/apt/lists/* && \
      apt-get autoremove -y && \
      apt-get clean all

COPY docker/templates/nginx.tmpl /etc/nginx/nginx.tmpl

COPY frontend "$FRONTEND_DIR"

RUN chown -R "$USER":"$GROUP" "$HOME"

USER "$USER"
WORKDIR "$HOME"

ENV PATH="$PREFIX/bin:/usr/local/node/bin:$PATH"
RUN \
      mkdir -p "$PREFIX/bin" && npm install -g yarn@"$YARN_VERSION" && \
      yarn global add angular-cli && cd "$FRONTEND_DIR" && yarn install && ng build --env="$BUILD" && \
      find . -maxdepth 1 ! \( -name 'dist' -o -name "." -o -name ".." \) -exec rm -rf {} + && \
      cd "$HOME" && \
      find . -maxdepth 1 ! \( -name "${FRONTEND_DIR##*/}" -o -name "." -o -name ".." \) -exec rm -rf {} +

EXPOSE 80 443

USER root
