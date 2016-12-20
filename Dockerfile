FROM alpine:3.4

ENV VIPS_VERSION 8.4.5

RUN apk add --update \
      ca-certificates \
      build-base \
      wget \
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
    && apk del --purge \
      build-base \
      wget \
    && rm -rf /var/cache/apk/*

ADD skrop /usr/local/bin/

ENTRYPOINT ["skrop"]
