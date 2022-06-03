FROM golang:1.18 as build

WORKDIR /go/src/app

COPY go.* ./
RUN go mod download

COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app -v .

FROM registry.access.redhat.com/ubi8/ubi-micro
COPY --from=build /go/src/app/app /
CMD ["/app"]
