# Use the official Go image as a parent image
FROM golang:1.21.3-bookworm

# Set the working directory in the container
WORKDIR /app

# Install necessary packages and tools
RUN apt-get update && apt-get install -y \
        wget \
        curl \
        gnupg \
        software-properties-common
        
# Copy the local code to the container's workspace
COPY . /app

# Install Go dependencies (if applicable)
RUN go mod download

# Run the application
CMD ./run_tests.sh


