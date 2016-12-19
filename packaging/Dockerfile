FROM golang:1.7.4-alpine

ENV VIPS_VERSION 8.4.5

RUN apk add --update \
      ca-certificates \
      build-base \
      wget \
      git \
      glib-dev \
      libxml2-dev \
    && wget http://www.vips.ecs.soton.ac.uk/supported/current/vips-${VIPS_VERSION}.tar.gz \
    && tar -zxvf vips-${VIPS_VERSION}.tar.gz \
    && cd vips-${VIPS_VERSION}/ \
    && ./configure \
      --disable-debug \
      --disable-static \
      --disable-introspection \
      --disable-dependency-tracking \
      --without-python \
      --without-orc \
      --without-fftw \
    && make \
    && make install \
    && cd ../ \
    && rm -rf vips-${VIPS_VERSION}/ \
    && rm vips-${VIPS_VERSION}.tar.gz \
    && rm -rf /var/cache/apk/*

ENV VIPSHOME /usr/local
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:$VIPSHOME/lib
ENV PATH $PATH:$VIPSHOME/bin
ENV PKG_CONFIG_PATH $PKG_CONFIG_PATH:$VIPSHOME/lib/pkgconfig

RUN go get github.com/tools/godep
