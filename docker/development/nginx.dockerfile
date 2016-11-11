FROM nginx:stable
MAINTAINER "Hugo Drumond" hugdru@gmail.com

# https://hub.docker.com/_/node/

ENV \
  USER="server" \
  GROUP="server" \
  NPM_CONFIG_LOGLEVEL="info" \
  NODE_VERSION="6.9.1" \
  SUEXEC_VERSION="v0.2" \
  SUEXEC_DOWNLOAD_SHA256="ec4acbd8cde6ceeb2be67eda1f46c709758af6db35cacbcde41baac349855e25" \
  WATCHMAN_VERSION="v4.7.0" \
  PERSISTENT_APT_PACKAGES="git libpcre++0 libpcre3 ca-certificates python make g++" \
  TEMPORARY_APT_PACKAGES="autoconf automake build-essential curl python-dev xz-utils libpcre3-dev libpcre++-dev wget ca-certificates" \
  DOCKERIZE_VERSION="v0.2.0" \
  DOCKERIZE_DOWNLOAD_SHA256="c0e2e33cfe066036941bf8f2598090bd8e01fdc05128490238b2a64cf988ecfb"


ENV HOME="/$USER"
ENV FRONTEND_DIR="$HOME/frontend"
ENV HTDOCS="$FRONTEND_DIR/dist"
ENV HOST_ENTRYPOINT_FILE="docker/development/entrypoints/nginx_entrypoint.sh"
ENV SCRIPTS_DIR="$HOME/bin"
ENV NPM_YARN_PACKAGES_DIR="$HOME/packages"

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
      git clone https://github.com/facebook/watchman.git /watchman && \
      cd watchman git checkout "$WATCHMAN_VERSION" && \
      ./autogen.sh && ./configure && make && make install && \
      mkdir /suexec && cd /suexec && \
      wget -O suexec.tar.gz "https://github.com/ncopa/su-exec/archive/$SUEXEC_VERSION.tar.gz" && \
      echo "$SUEXEC_DOWNLOAD_SHA256  suexec.tar.gz" | sha256sum -c - && \
      tar -xvzf suexec.tar.gz --strip-components 1 && make && mv su-exec /usr/local/bin && cd .. && rm -rf /suexec && \
      curl -fsSL "https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz" -o dockerize.tar.gz && \
      echo "$DOCKERIZE_DOWNLOAD_SHA256  dockerize.tar.gz" | sha256sum -c - && \
      tar -C /usr/local/bin -xzvf dockerize.tar.gz && rm dockerize.tar.gz && \
      apt-get remove --purge -y $TEMPORARY_APT_PACKAGES && \
      rm -rf /var/lib/apt/lists/* /watchman && \
      apt-get autoremove -y && \
      apt-get clean all

VOLUME "$FRONTEND_DIR"
COPY docker/templates/nginx.tmpl /etc/nginx/nginx.tmpl

WORKDIR "$SCRIPTS_DIR"
COPY "$HOST_ENTRYPOINT_FILE" entrypoint.sh

RUN chown -R "$USER":"$GROUP" "$HOME" && chmod +x "$SCRIPTS_DIR/entrypoint.sh"

USER "$USER"
WORKDIR "$HOME"

ENV PREFIX="$NPM_YARN_PACKAGES_DIR" YARN_PREFIX="$NPM_YARN_PACKAGES_DIR"
ENV PATH="$NPM_YARN_PACKAGES_DIR/bin:$SCRIPTS_DIR:/usr/local/node/bin:$PATH"
RUN mkdir -p "$NPM_YARN_PACKAGES_DIR/bin" && npm install -g yarn && yarn global add angular-cli

EXPOSE 80 443 4200

USER root
