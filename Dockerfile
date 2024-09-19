
FROM --platform=linux/amd64 golang:1.19-alpine as builder

RUN mkdir -p /app/api
WORKDIR /app/api

ENV GOOS=linux CGO_ENABLED=1

RUN set -ex && \
  apk add --no-progress --no-cache \
  gcc \
  musl-dev 
COPY . .
RUN go install --tags musl ./...
EXPOSE 8085

CMD [ "api" ]