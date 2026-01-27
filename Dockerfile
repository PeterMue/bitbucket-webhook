# Build stage
FROM --platform=$BUILDPLATFORM golang:latest AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o server .

# Final stage
FROM alpine:latest
COPY --from=build /app/server /server
ENTRYPOINT ["/server"]
