FROM golang:1.25-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY api-go/ ./api-go/
RUN CGO_ENABLED=0 go build -trimpath -o /api ./api-go/cmd/api

FROM scratch
COPY --from=build /api /api
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/api"]