# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
# 默认走官方代理；国内若 `go mod download` 超时，构建时传入：
# --build-arg GOPROXY=https://goproxy.cn,direct --build-arg GOSUMDB=sum.golang.google.cn
ARG GOPROXY=https://proxy.golang.org,direct
ARG GOSUMDB=sum.golang.org
ENV GOPROXY=${GOPROXY} GOSUMDB=${GOSUMDB}
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/learn-api ./cmd/learn-api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates \
	&& adduser -D -u 65532 appuser
WORKDIR /app
COPY --from=build /out/learn-api /app/learn-api
USER appuser
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/app/learn-api"]
