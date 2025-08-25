# Frontend build stage
FROM node:22.18.0-alpine3.22 AS frontend-build-env

WORKDIR /app

COPY jmeter/site/ ./

RUN npm install && npm run build

# Backend build stage
FROM golang:1.25.0-alpine3.22 AS backend-build-env

ENV PROJECT_GROUP=derek82511
ENV PROJECT_NAME=jt

RUN apk update && apk add --no-cache gcc libc-dev

ADD go.mod /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/
ADD go.sum /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/
ADD main.go /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/
ADD config /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/config/
ADD model /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/model/
ADD service /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/service/

WORKDIR /go/src/${PROJECT_GROUP}/${PROJECT_NAME}

RUN go mod init ${PROJECT_GROUP}/${PROJECT_NAME} || echo "go.mod already exists" \
  && go mod tidy && go build -o bin/app

WORKDIR /source

RUN wget https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.8%2B9/OpenJDK21U-jdk_x64_alpine-linux_hotspot_21.0.8_9.tar.gz \
  && wget https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-5.6.3.tgz \
  && tar -xzvf OpenJDK21U-jdk_x64_alpine-linux_hotspot_21.0.8_9.tar.gz \
  && tar -xzvf apache-jmeter-5.6.3.tgz \
  && mv jdk-21.0.8+9 java \
  && mv apache-jmeter-5.6.3 jmeter

# Runtime build stage
FROM alpine:3.22.1

ENV PROJECT_GROUP=derek82511
ENV PROJECT_NAME=jt

RUN mkdir /usr/local/java \
  && mkdir /usr/local/jmeter

COPY --from=backend-build-env /source/java /usr/local/java/
COPY --from=backend-build-env /source/jmeter /usr/local/jmeter/

RUN echo "client.rmi.localport=10090" | tee -a /usr/local/jmeter/bin/jmeter.properties \
  && echo "server.rmi.ssl.disable=true" | tee -a /usr/local/jmeter/bin/jmeter.properties \
  && echo "java.rmi.server.hostname=SERVER_HOST" | tee -a /usr/local/jmeter/bin/system.properties \
  && echo "jmeter.reportgenerator.graph.responseTimeVsRequest.exclude_controllers=false" | tee -a /usr/local/jmeter/bin/user.properties \
  && echo "jmeter.reportgenerator.graph.latencyVsRequest.exclude_controllers=false" | tee -a /usr/local/jmeter/bin/user.properties \
  && echo "export JAVA_HOME=/usr/local/java" | tee -a /etc/profile \
  && echo "export JRE_HOME=\${JAVA_HOME}/jre" | tee -a /etc/profile \
  && echo "export CLASSPATH=.:\${JAVA_HOME}/lib:\${JRE_HOME}/lib" | tee -a /etc/profile \
  && echo "export PATH=\${JAVA_HOME}/bin:\$PATH" | tee -a /etc/profile \
  && echo "export PATH=/usr/local/jmeter/bin:\$PATH" | tee -a /etc/profile

RUN mkdir /root/jmeter \
  && mkdir /root/jmeter/shell \
  && mkdir /root/jmeter/scripts \
  && mkdir /root/jmeter/reports \
  && mkdir /root/jmeter/logs \
  && mkdir /root/jmeter/data

ADD ./jmeter/shell/run.sh /root/jmeter/shell/
ADD ./jmeter/shell/recovery.sh /root/jmeter/shell/

RUN chmod +x /root/jmeter/shell/run.sh \
  && chmod +x /root/jmeter/shell/recovery.sh

COPY --from=frontend-build-env /app/dist /root/jmeter/site/dist/

COPY --from=backend-build-env /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/bin/app /runtime/

ADD ./env/server.sh /

RUN chmod +x /server.sh

EXPOSE 8088 10090 10091 10092

ENV SERVER_HOST=127.0.0.1

ENTRYPOINT ["/server.sh"]
