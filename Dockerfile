# We specify the base image we need for our
# go application
FROM golang:1.12.6-alpine3.10
# We update and upgrade packages.
# We install git
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
# Set timezone
RUN apk add --no-cache tzdata
ENV TZ Asia/Almaty
# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app
# We copy everything in the root directory
# into our /app directory
ADD . /app
# We specify that we now wish to execute
# any further commands inside our /app
# directory
WORKDIR /app
# go get modules
RUN go get github.com/ivahaev/russian-time
RUN go get github.com/jasonlvhit/gocron
RUN go get github.com/joho/godotenv
RUN go get github.com/getsentry/raven-go
# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .
# Our start command which kicks off
# our newly created binary executable
CMD ["/app/main"]