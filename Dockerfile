FROM golang:1.20

WORKDIR /app

COPY go.mod ./
COPY app ./app
COPY club ./club
COPY entities ./entities
COPY utils ./utils
COPY main.go ./


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /computer-club

# Run
CMD [ "/computer-club", "input.txt"]
