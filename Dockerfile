# From image that contains go
FROM golang:1.9

# Create the directory
RUN mkdir -p /go/src/github.com/jsotogaviard/discount

# Copy all the relevant information
ADD . /go/src/github.com/jsotogaviard/discount

# Configure as current directory
WORKDIR /go/src/github.com/jsotogaviard/discount

# Get and build swagger
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
RUN swagger generate server -f ./api/discount.yaml
RUN swagger generate client -f ./api/discount.yaml

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# Build the binary file
WORKDIR /go/src/github.com/jsotogaviard/discount/cmd/checkout-server
RUN go build -o main .

# Launch the main
CMD ["/go/src/github.com/jsotogaviard/discount/cmd/checkout-server/main"]

EXPOSE 8000