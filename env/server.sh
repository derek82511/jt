#!/bin/sh

sed -i "s/java.rmi.server.hostname=SERVER_HOST/java.rmi.server.hostname=$SERVER_HOST/" /usr/local/jmeter/bin/system.properties

source /etc/profile

/runtime/app
