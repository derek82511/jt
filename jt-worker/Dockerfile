# Build stage to download assets
FROM alpine:3.22.1 AS build-env

WORKDIR /source

RUN wget https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.8%2B9/OpenJDK21U-jdk_x64_alpine-linux_hotspot_21.0.8_9.tar.gz \
  && wget https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-5.6.3.tgz \
  && tar -xzvf OpenJDK21U-jdk_x64_alpine-linux_hotspot_21.0.8_9.tar.gz \
  && tar -xzvf apache-jmeter-5.6.3.tgz \
  && mv jdk-21.0.8+9 java \
  && mv apache-jmeter-5.6.3 jmeter

# Runtime build stage
FROM alpine:3.22.1

RUN mkdir -p /usr/local/java \
  && mkdir -p /usr/local/jmeter

COPY --from=build-env /source/java /usr/local/java/
COPY --from=build-env /source/jmeter /usr/local/jmeter/

RUN echo "server_port=1099" | tee -a /usr/local/jmeter/bin/jmeter.properties \
  && echo "server.rmi.localport=1099" | tee -a /usr/local/jmeter/bin/jmeter.properties \
  && echo "server.rmi.ssl.disable=true" | tee -a /usr/local/jmeter/bin/jmeter.properties \
  && echo "java.rmi.server.hostname=SERVER_HOST" | tee -a /usr/local/jmeter/bin/system.properties \
  && echo "export JAVA_HOME=/usr/local/java" | tee -a /etc/profile \
  && echo "export JRE_HOME=\${JAVA_HOME}/jre" | tee -a /etc/profile \
  && echo "export CLASSPATH=.:\${JAVA_HOME}/lib:\${JRE_HOME}/lib" | tee -a /etc/profile \
  && echo "export PATH=\${JAVA_HOME}/bin:\$PATH" | tee -a /etc/profile \
  && echo "export PATH=/usr/local/jmeter/bin:\$PATH" | tee -a /etc/profile

ADD ./env/server.sh /

RUN chmod +x /server.sh

EXPOSE 1099

ENV SERVER_HOST=127.0.0.1

ENTRYPOINT ["/server.sh"]
