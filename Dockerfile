ARG GO_VERSION=1.20
ARG ALPINE_VERSION=3.16

# Builder 
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV PROJECT=tf-plan-analyzer

WORKDIR ${PROJECT}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -a -o /${PROJECT} .

### Scratch with build in docker
FROM scratch as goreleaser
COPY --from=builder /tf-plan-analyzer /bin/
ENTRYPOINT ["/bin/tf-plan-analyzer"]
