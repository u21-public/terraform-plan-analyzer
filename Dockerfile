### Scratch with build in docker
FROM golang:alpine as build
# Redundant, current golang images already include ca-certificates
RUN apk --no-cache add ca-certificates

FROM scratch as goreleaser
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY terraform-plan-analyzer /terraform-plan-analyzer
ENTRYPOINT ["/terraform-plan-analyzer"]
