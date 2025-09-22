FROM golang:tip-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o media_updater ./scripts/media_updater.go

ENTRYPOINT ["./media_updater"]