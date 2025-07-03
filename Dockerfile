FROM golang:1.24
LABEL authors="elnerribeiro"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /gopoc
RUN chmod +x exec.sh
EXPOSE 8000

# Run
CMD ["./exec.sh"]

