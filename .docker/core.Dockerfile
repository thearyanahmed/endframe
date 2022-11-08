FROM golang:1.19

ENV APP_DIR /app

# Install CompileDaemon for auto installing binary on changes
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY cmd/core/ core/ .env Makefile $APP_DIR
WORKDIR $APP_DIR

RUN make deps-core