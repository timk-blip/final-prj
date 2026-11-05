FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o todo ./cmd/todo

FROM ubuntu:latest
WORKDIR /app
COPY --from=build /src/todo /app/todo
COPY web /app/web
ENV TODO_PORT=7540
EXPOSE 7540
CMD ["/app/todo"]