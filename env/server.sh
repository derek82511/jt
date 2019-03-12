#!/bin/bash

sed -i "s/java.rmi.server.hostname=SERVER_HOST/java.rmi.server.hostname=$SERVER_HOST/" /usr/local/jmeter/apache-jmeter-5.1/bin/system.properties

source /etc/profile

/runtime/app
