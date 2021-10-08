FROM golang:alpine as builder
WORKDIR /ghb
RUN apk add make
COPY . /ghb
RUN make

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/mfinelli/github-backup
RUN apk add --no-cache git gnupg gzip openssh-client-default tar
COPY --from=builder /ghb/ghb /usr/bin/ghb
ENTRYPOINT ["ghb"]
