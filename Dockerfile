FROM golang:1.17

ENV DEVZONE_PORT=8082
WORKDIR /usr/src/devZoneDeployment

ADD go.mod go.sum ./
RUN go mod download && go mod verify

ADD . .
ADD config.example.yaml /etc/devZoneDeployment/config.default.yaml
RUN go build -v -o /usr/local/bin/devZoneDeployment

EXPOSE $DEVZONE_PORT
VOLUME /etc/devZoneDeployment/
VOLUME /etc/openvpn/server/

CMD ["devZoneDeployment", "--config", "/etc/devZoneDeployment/config.yaml"]