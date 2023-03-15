FROM timescale/timescaledb:latest-pg15

# Set the environment variables
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD UrWPqf5kpG5Rpq3r!
ENV POSTGRES_DB homework
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432

# Expose the default PostgreSQL port
EXPOSE 5432

# Install go
RUN apk update && apk add --no-cache git go

# Set the working directory to /app
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Install the Go modules
RUN go mod download

# Copy the Go application files
COPY . .

# Build the Go application
RUN go build -o tsb
