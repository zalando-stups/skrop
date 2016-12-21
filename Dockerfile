FROM alpine:3.4

ARG ROUTES_FILE

ADD packaging/install-vips.sh install-vips.sh
RUN apk add --update \
    bash \
  && ./install-vips.sh \
  && apk del --purge \
    build-base \
    wget \
    bash \
  && rm -rf /var/cache/apk/* \
  && rm install-vips.sh

ADD skrop /usr/local/bin/
ADD $ROUTES_FILE skrop.eskip

ENTRYPOINT ["skrop", "-routes-file", "skrop.eskip"]
CMD $SKROP_ARGS
