# What is JT?

JT (JMeter Tool) is a web application which packages the Apache JMeter open source software.

JT makes it very convenience to run jmeter script. You can easily upload your jmeter script(.jmx), run it, and check the reports on the web.

This image use apache-jmeter-5.1 and openjdk-8.

Docker Hub: [https://hub.docker.com/r/derek82511/jt](https://hub.docker.com/r/derek82511/jt)

GitHub: [https://github.com/derek82511/jt](https://github.com/derek82511/jt)

# How to use this image

## Prepare

Create a directory named jmeter includes 4 path (scripts, reports, logs, data).

```
- jmeter
    - scripts  (uploaded jmx scripts)
    - reports  (generated reports)
    - logs     (web, application and sql log)
    - data     (application data)
```

## Simple usage

```console
docker run -it \
  -p 10080:10080 \
  -e "TZ=Asia/Taipei" \
  -v $PWD/jmeter/scripts:/root/jmeter/scripts \
  -v $PWD/jmeter/reports:/root/jmeter/reports \
  -v $PWD/jmeter/logs:/root/jmeter/logs \
  -v $PWD/jmeter/data:/root/jmeter/data \
  --name jt \
  derek82511/jt:1.2.0
```

Open `http://localhost:10080` or `http://host-ip:10080` in your browser.

## Using host networking

Sometimes, use host network would improve your testing performance if your host has below setting.
```
net.ipv4.tcp_fin_timeout=30  
net.ipv4.tcp_tw_reuse=1
net.ipv4.tcp_tw_recycle=1
```

```console
docker run -it \
  --network=host \
  -e "TZ=Asia/Taipei" \
  -v $PWD/jmeter/scripts:/root/jmeter/scripts \
  -v $PWD/jmeter/reports:/root/jmeter/reports \
  -v $PWD/jmeter/logs:/root/jmeter/logs \
  -v $PWD/jmeter/data:/root/jmeter/data \
  --name jt \
  derek82511/jt:1.2.0
```

Open `http://localhost:10080` or `http://host-ip:10080` in your browser.

## Remote Testing

Support since v1.1.0.

You must prepare your JMeter worker server and start it. (Or you can use [jt-worker](https://hub.docker.com/r/derek82511/jt-worker) in your worker server)
* Running exactly the same version of JMeter (apache-jmeter-5.0).
* Using the same version of Java on all systems (openjdk-8). Using different versions of Java may work but is discouraged.

Following configuration is required in jmeter.properties file:

```
server.rmi.ssl.disable=true
```

```console
docker run -it \
  -p 10080:10080 \
  -p 10090:10090 \
  -p 10091:10091 \
  -p 10092:10092 \
  -e "TZ=Asia/Taipei" \
  -e "SERVER_HOST=[YOUR HOST IP]" \
  -v $PWD/jmeter/scripts:/root/jmeter/scripts \
  -v $PWD/jmeter/reports:/root/jmeter/reports \
  -v $PWD/jmeter/logs:/root/jmeter/logs \
  -v $PWD/jmeter/data:/root/jmeter/data \
  --name jt \
  derek82511/jt:1.2.0
```

or

```console
docker run -it \
  --network=host \
  -e "TZ=Asia/Taipei" \
  -e "SERVER_HOST=[YOUR HOST IP]" \
  -v $PWD/jmeter/scripts:/root/jmeter/scripts \
  -v $PWD/jmeter/reports:/root/jmeter/reports \
  -v $PWD/jmeter/logs:/root/jmeter/logs \
  -v $PWD/jmeter/data:/root/jmeter/data \
  --name jt \
  derek82511/jt:1.2.0
```

Then, you can assign the worker server's ip:port when creating job with remote mode.
