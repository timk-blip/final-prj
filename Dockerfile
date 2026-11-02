FROM ubuntu:latest
WORKDIR /app
COPY todo-arm64 /app/todo
COPY web /app/web
ENV TODO_PORT=7540
EXPOSE 7540
CMD ["/app/todo"]