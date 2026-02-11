# build stage
FROM golang:1.22.2-alpine AS build
WORKDIR /app

# Print every command and output
RUN set -x

# cache dependencies
COPY go.mod go.sum ./
# RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

# final stage
FROM alpine:latest
RUN set -x
RUN adduser -D app
WORKDIR /home/app

# Copy the server binary
COPY --from=build /app/server .

# Create static directory and copy static assets
RUN mkdir -p static
COPY --from=build /app/garnetstar.ico ./static/
COPY --from=build /app/garnetstar.jpeg ./static/

USER app
EXPOSE 8080
ENTRYPOINT ["./server"]