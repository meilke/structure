FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build github.com/meilke/structure 
CMD ["/app/structure"]
