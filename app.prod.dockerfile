# Start from golang base image
FROM golang:latest as builder

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/linkwind/app
WORKDIR /go/src/linkwind/app

COPY .env ./
COPY .env.dev ./

# Copy go mod and sum files 
COPY ./app/go.mod ./app/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /go/src/linkwind/app/main .
COPY --from=builder /go/src/linkwind/app/.env .
COPY --from=builder /go/src/linkwind/app/templates ./templates
COPY --from=builder /go/src/linkwind/app/public ./public
COPY --from=builder /go/src/linkwind/app/data/sql_scripts ./sql_scripts

# Expose port to the outside world
EXPOSE 8080

CMD ["./main"];

# if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
# if production setting will build binary
# CMD if [ ${APP_ENV} = production ]; \
#     then \
#     ["./main"]; \
#     else \
#     go get github.com/pilu/fresh && \
#     fresh; \
#     fi