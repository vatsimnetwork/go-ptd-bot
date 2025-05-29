FROM golang:1.23 as build

WORKDIR /workspace
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /workspace/bin/bot /workspace/cmd/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=build /workspace/bin/bot /
USER 65532:65532

ENTRYPOINT ["/bot"]
