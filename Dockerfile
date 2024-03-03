# https://docs.docker.com/language/golang/build-images/#multi-stage-builds

# Build the application from source
FROM golang:1.21.4-alpine3.18 AS build-stage
# FROM golang:1.21.4-bookworm AS build-stage

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

# COPY *.go ./
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-s -w" -o /app ./cmd/main/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app ./cmd/main/main.go

# # Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
# https://github.com/GoogleContainerTools/distroless/tree/main/base
FROM gcr.io/distroless/static-debian12:nonroot AS build-release-stage

COPY --from=build-stage /app /rinha-2024q1-crebito

EXPOSE 3000

ENTRYPOINT ["/rinha-2024q1-crebito"]