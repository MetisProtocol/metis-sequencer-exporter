# syntax=docker/dockerfile:1
FROM golang:1.21.6 as builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go install .

FROM gcr.io/distroless/base-debian12
COPY --from=builder /go/bin/metis-sequencer-exporter /usr/local/bin/
ENTRYPOINT [ "metis-sequencer-exporter" ]
