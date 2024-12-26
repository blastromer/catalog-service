# Step 1: Use the official Golang image as a base image
FROM golang:1.20-alpine AS builder

# Step 2: Set the working directory in the container
WORKDIR /app

# Step 3: Copy the Go module files
COPY go.mod go.sum ./

# Step 4: Download Go dependencies
RUN go mod tidy

# Step 5: Copy the source code into the container
COPY . .

# Step 6: Build the Go binary
RUN go build -o main .

# Step 7: Create a new, smaller image for production
FROM alpine:latest

# Step 8: Install required libraries (if needed)
RUN apk --no-cache add ca-certificates

# Step 9: Set the working directory in the production image
WORKDIR /root/

# Step 10: Copy the binary from the builder stage to the production stage
COPY --from=builder /app/main .

# Step 11: Expose the port that the Go app will run on
EXPOSE 8080

# Step 12: Command to run the Go application
CMD ["./main"]
