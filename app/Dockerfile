FROM golang:1.13
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/golang/mock/gomock
RUN go get github.com/golang/mock/mockgen
RUN mkdir $GOPATH/src/app
ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN mockgen -destination=mocks/mock_httpclient.go -package=mocks -source=prometheus_collector.go HttpClient
RUN go test .
RUN go build -o prometheus_collector .

CMD ["./prometheus_collector"]