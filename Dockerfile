# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26 AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o server .

# Final stage
FROM alpine:3.23

RUN apk update \
     && apk upgrade \
     && apk add --no-cache --update ca-certificates \
     && addgroup -S webhook \
     && adduser -S webhook -G webhook

WORKDIR /home/webhook

COPY --from=build /app/server /bin/bitbucket-webhook

EXPOSE 3000

ENTRYPOINT ["/bin/bitbucket-webhook"]
