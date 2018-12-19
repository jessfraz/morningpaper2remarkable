FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/jessfraz/morningpaper2remarkable

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/jessfraz/morningpaper2remarkable \
	&& make static \
	&& mv morningpaper2remarkable /usr/bin/morningpaper2remarkable \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/morningpaper2remarkable /usr/bin/morningpaper2remarkable
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

RUN adduser -D -u 1000 user \
  && chown -R user /home/user

USER user

ENV USER user

WORKDIR /home/user

ENTRYPOINT [ "morningpaper2remarkable" ]
CMD [ "--help" ]
