### Scratch with build in docker
FROM scratch as goreleaser
COPY terraform-plan-analyzer /terraform-plan-analyzer
ENTRYPOINT ["/terraform-plan-analyzer"]
