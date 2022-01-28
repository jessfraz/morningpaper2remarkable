FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/pseudo-su/morningpaper2remarkable

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/pseudo-su/morningpaper2remarkable \
	&& make static \
	&& mv morningpaper2remarkable /usr/bin/morningpaper2remarkable \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/morningpaper2remarkable /usr/bin/morningpaper2remarkable
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "morningpaper2remarkable" ]
CMD [ "--help" ]
