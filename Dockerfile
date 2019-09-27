# build stage
FROM golang:stretch AS build-env

ENV PROJECT_GROUP derek82511
ENV PROJECT_NAME jt

RUN apt-get update && apt-get install -y --no-install-recommends git \
  && go get -u github.com/kardianos/govendor

ADD ./main.go /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/
ADD ./config /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/config/
ADD ./model /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/model/
ADD ./service /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/service/
ADD ./vendor/vendor.json /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/vendor/

WORKDIR /go/src/${PROJECT_GROUP}/${PROJECT_NAME}
RUN govendor sync \
  && go build -o bin/app \
  && wget https://download.java.net/openjdk/jdk8u40/ri/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz \
  && wget http://ftp.tc.edu.tw/pub/Apache//jmeter/binaries/apache-jmeter-5.1.tgz

# final stage
FROM debian:stretch-slim

ENV PROJECT_GROUP derek82511
ENV PROJECT_NAME jt

RUN mkdir /usr/local/java \
  && mkdir /usr/local/jmeter

COPY --from=build-env /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz /usr/local/java
COPY --from=build-env /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/apache-jmeter-5.1.tgz /usr/local/jmeter

RUN tar -xzvf /usr/local/java/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz -C /usr/local/java \
  && tar -xzvf /usr/local/jmeter/apache-jmeter-5.1.tgz -C /usr/local/jmeter \
  && echo "client.rmi.localport=10090" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/jmeter.properties \
  && echo "server.rmi.ssl.disable=true" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/jmeter.properties \
  && echo "java.rmi.server.hostname=SERVER_HOST" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/system.properties \
  && echo "jmeter.reportgenerator.graph.responseTimeVsRequest.exclude_controllers=false" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/user.properties \
  && echo "jmeter.reportgenerator.graph.latencyVsRequest.exclude_controllers=false" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/user.properties \
  && echo "export JAVA_HOME=/usr/local/java/java-se-8u40-ri" | tee -a /etc/profile \
  && echo "export JRE_HOME=\${JAVA_HOME}/jre" | tee -a /etc/profile \
  && echo "export CLASSPATH=.:\${JAVA_HOME}/lib:\${JRE_HOME}/lib" | tee -a /etc/profile \
  && echo "export PATH=\${JAVA_HOME}/bin:\$PATH" | tee -a /etc/profile \
  && echo "export PATH=/usr/local/jmeter/apache-jmeter-5.1/bin:\$PATH" | tee -a /etc/profile

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

ADD ./jmeter/site/dist /root/jmeter/site/dist/

COPY --from=build-env /go/src/${PROJECT_GROUP}/${PROJECT_NAME}/bin/app /runtime/

ADD ./env/server.sh /

RUN chmod +x /server.sh

EXPOSE 10080 10090 10091 10092

ENV SERVER_HOST 127.0.0.1

ENTRYPOINT ["/server.sh"]
