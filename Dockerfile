FROM golang:1.22

WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
