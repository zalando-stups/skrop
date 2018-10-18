FROM danpersa/alpine-vips:8.6.5

ARG ROUTES_FILE

ADD skrop /usr/local/bin/
ADD $ROUTES_FILE skrop.eskip

ENTRYPOINT skrop -routes-file skrop.eskip ${SKROP_ARGS}
