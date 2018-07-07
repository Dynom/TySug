FROM golang:1.10 as build

WORKDIR /go/src/github.com/Dynom/TySug
COPY . .

RUN go test -test.short -test.v -test.race ./...
RUN CGO_ENABLED=0 go install -v -a -ldflags "-w" ./...

FROM gcr.io/distroless/base as base

FROM scratch

LABEL description="The TySug webservice Docker image. Suggesting typo-alternatives"

COPY --from=base ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=base ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build /go/bin/TySug /
COPY ["config.yml", "/"]

CMD ["/TySug"]