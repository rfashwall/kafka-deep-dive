FROM golang:1.17.6-alpine AS build

# WORKDIR /src/
# COPY main.go go.* /src/
# RUN CGO_ENABLED=0 go build -o inventory cmd/inventory

# FROM scratch
# COPY --from=build /bin/inventory /bin/inventory
# ENTRYPOINT ["/bin/inventory"]

# Set the current working directory inside the container

# Set the current working directory inside the container
WORKDIR /build
# RUN git config --global url
#Pass the content of the private key into the container
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
#Github requires a private key with strict permission settings
RUN chmod 600 /root/.ssh/id_rsa
#Add Github to known hosts
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

 #Copy go.mod, go.sum files and download deps
 COPY go.mod go.sum ./
 RUN go mod init

# Copy sources to the working directory
COPY . .

# Build the Go app
ARG project
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o server $project

# Start a new stage from busybox
FROM busybox:latest

WORKDIR /dist

# Copy the build artifacts from the previous stage
COPY --from=builder /build/server .

# Run the executable
CMD ["./server"]