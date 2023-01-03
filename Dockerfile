FROM golang:1.19

WORKDIR /mpromt

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o mpromt .

ENV ADDR=0.0.0.0
EXPOSE 2112

CMD [ "./mpromt" ]