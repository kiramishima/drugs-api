FROM golang:1.22-alpine as build
ENV API_NAME=api_drugs
ENV APP_SECRET=B@nk4I
  # Database
ENV DATABASE_DRIVER=pgx
ENV DATABASE_URL=postgres://postgres:postgres@database:5432/ionix
ENV DATABASE_MAX_OPEN_CONNECTIONS=25
ENV DATABASE_MAX_IDDLE_CONNECTIONS=25
ENV DATABASE_MAX_IDDLE_TIME=15m
  # HTTP
ENV HTTP_SERVER_IDLE_TIMEOUT=60s
ENV PORT: 8080
ENV HTTP_SERVER_READ_TIMEOUT=1s
ENV HTTP_SERVER_WRITE_TIMEOUT=2s
# JWT
ENV JWT_PRIVATE_KEY=SecretMedicament
ENV TOKEN_TTL=300
# Context
ENV CONTEXT_TIMEOUT=10

RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
ENV GOPROXY https://proxy.golang.org,direct
RUN go mod download
ENV CGO_ENABLED=0
RUN GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o $API_NAME ./cmd/main.go

FROM scratch as serve
WORKDIR /app
COPY --from=build /app/$API_NAME .
CMD ["/app/api_drugs"]