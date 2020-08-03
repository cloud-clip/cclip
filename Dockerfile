# Go 1.14.x based on Alpine Linux
FROM golang:1.14-alpine

# make and define working directory
RUN mkdir /app
ADD . /app
WORKDIR /app

# copy all files
COPY . .

# now build ...
RUN go build -o cclip .

# ... and run
CMD ["./cclip"]
