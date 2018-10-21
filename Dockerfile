# STEP 1 build executable binary
FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/stewartmatheson/greg
WORKDIR $GOPATH/src/github.com/stewartmatheson/greg
#get dependancies
#you can also use dep
RUN apk add git
RUN go get -d -v
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags=”-w -s” -o /go/bin/greg
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/greg .
# STEP 2 build a small image
# start from scratch
FROM scratch
# Copy our static executable
COPY --from=builder /go/bin/greg /go/bin/greg
ENTRYPOINT ["/go/bin/greg"]
