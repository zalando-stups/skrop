FROM alpine:3.8

# Build-time metadata as defined at http://label-schema.org
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Alpine Mozjpeg Vips" \
      org.label-schema.description="Alpine with mozjpeg and vips installed" \
      org.label-schema.url="https://github.com/zalando-stups/skrop" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/zalando-stups/skrop" \
      org.label-schema.version="3.3.1-8.7.0" \
      org.label-schema.schema-version="1.0"

ADD docker/install-vips.sh install-vips.sh

RUN apk add --update \
    bash \
  && ./install-vips.sh \
  && apk del --purge \
    build-base \
    wget \
    bash \
  && rm -rf /var/cache/apk/* \
  && rm install-vips.sh
