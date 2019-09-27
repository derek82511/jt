# build stage
FROM golang:stretch AS build-env

WORKDIR /src

RUN wget https://download.java.net/openjdk/jdk8u40/ri/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz \
  && wget http://ftp.tc.edu.tw/pub/Apache//jmeter/binaries/apache-jmeter-5.1.tgz

# final stage
FROM debian:stretch-slim

RUN mkdir /usr/local/java \
  && mkdir /usr/local/jmeter

COPY --from=build-env /src/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz /usr/local/java
COPY --from=build-env /src/apache-jmeter-5.1.tgz /usr/local/jmeter

RUN tar -xzvf /usr/local/java/openjdk-8u40-b25-linux-x64-10_feb_2015.tar.gz -C /usr/local/java \
  && tar -xzvf /usr/local/jmeter/apache-jmeter-5.1.tgz -C /usr/local/jmeter \
  && echo "server_port=1099" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/jmeter.properties \
  && echo "server.rmi.localport=1099" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/jmeter.properties \
  && echo "server.rmi.ssl.disable=true" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/jmeter.properties \
  && echo "java.rmi.server.hostname=SERVER_HOST" | tee -a /usr/local/jmeter/apache-jmeter-5.1/bin/system.properties \
  && echo "export JAVA_HOME=/usr/local/java/java-se-8u40-ri" | tee -a /etc/profile \
  && echo "export JRE_HOME=\${JAVA_HOME}/jre" | tee -a /etc/profile \
  && echo "export CLASSPATH=.:\${JAVA_HOME}/lib:\${JRE_HOME}/lib" | tee -a /etc/profile \
  && echo "export PATH=\${JAVA_HOME}/bin:\$PATH" | tee -a /etc/profile \
  && echo "export PATH=/usr/local/jmeter/apache-jmeter-5.1/bin:\$PATH" | tee -a /etc/profile

ADD ./env/server.sh /

RUN chmod +x /server.sh

EXPOSE 1099

ENV SERVER_HOST 127.0.0.1

ENTRYPOINT ["/server.sh"]
