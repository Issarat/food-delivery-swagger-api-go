FROM golang:1.25-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY food-delivery-orders-api-go/cmd ./food-delivery-orders-api-go/cmd
COPY food-delivery-orders-api-go/internal ./food-delivery-orders-api-go/internal
COPY food-delivery-orders-api-go/docs ./food-delivery-orders-api-go/docs
RUN CGO_ENABLED=0 go build -trimpath -o /api ./food-delivery-orders-api-go/cmd/api

FROM scratch
COPY --from=build /api /api
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/api"]