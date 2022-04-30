FROM fedora:latest

LABEL maintainer="Ricardo Ferreira <riferrei@riferrei.com>"

# Install pre-reqs
RUN dnf install wget -y
RUN dnf install gcc -y

# Install librdkafka
RUN rpm --import https://packages.confluent.io/rpm/7.1/archive.key
COPY confluent.repo /etc/yum.repos.d
RUN dnf install librdkafka-devel -y

# Install Go v1.18
RUN wget https://dl.google.com/go/go1.18.linux-amd64.tar.gz && tar -xvf go1.18.linux-amd64.tar.gz && rm go1.18.linux-amd64.tar.gz
RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH="${GOROOT}/bin:${PATH}"

# Build the producer
WORKDIR /app
COPY go.mod .
COPY app.go .
RUN go build -o app .
RUN rm app.go && rm go.*

CMD ["./app"]