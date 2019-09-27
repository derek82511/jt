# How to use this image

## Simple usage

```console
docker run -it \
  -p 1099:1099 \
  -e "TZ=Asia/Taipei" \
  -e "SERVER_HOST=[YOUR HOST IP]" \
  --name jt-worker \
  derek82511/jt-worker:1.3.0
```

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
  -e "SERVER_HOST=[YOUR HOST IP]" \
  --name jt-worker \
  derek82511/jt-worker:1.3.0
```

## Reference image

[jt](https://hub.docker.com/r/derek82511/jt)
