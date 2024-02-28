# syntax=docker/dockerfile:1
FROM golang:1.22.0 
#RUN adduser -D -u 1000 -g 1000 henry
# Set destination for COPY
WORKDIR /app
#ENV LOG_DIR /app/logs
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy

COPY cmd/*.go ./cmd/
COPY controllers/*.go ./controllers/
COPY internal/middlewares/*.go ./internal/middlewares/
COPY internal/models/*.go ./internal/models/
COPY internal/utils/*.go ./internal/utils/
COPY logs/* ./logs/
COPY templates/*.html ./templates/
COPY ./assets/* ./assets/
COPY ./assets/css/* ./assets/css/
VOLUME /app/logs
#USER henry
# /var/lib/docker/volumes/logs/ à l'extérieur du container
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o 118_session_ok cmd/main.go


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080
# Run
CMD ["./118_session_ok"]
# Pour construire le container avec le nom 118_session_ok
# docker build -t 118_session_ok .
# Pour lancer l'application 
# docker run -it --rm -p 8080:8080 -v /home/henry/go/src/118_session_ok/logs:/app/logs 118_session_ok