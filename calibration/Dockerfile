FROM golang:1.18 as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o ./bin/calib

FROM gcr.io/nsx-sm/photon:3.0
WORKDIR /bin
COPY --from=builder /app/bin/calib .
USER 65532:65532
ENTRYPOINT ["/bin/calib"]
