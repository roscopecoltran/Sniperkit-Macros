Usage
=====

```bash
$ docker run -d --name mongo mongo
$ docker run -d --name elasticsearch blacktop/elasticsearch
$ docker run -d --name postgres -e POSTGRES_PASSWORD=cuckoo postgres
# Start cuckoo API
$ docker run -d --name cuckoo-api \
				--link postgres \
				-p 8000:1337 \
				blacktop/cuckoo:2.0 api
# Start cuckoo web UI				
$ docker run -d --name cuckoo-web \
				--link mongo \
				--link elasticsearch \
				-p 80:31337 \
				blacktop/cuckoo:2.0 web
```

> **NOTE:** If you want to customize the cuckoo configuration before launching you can link the **conf** folder into the container like so: `docker run -d -v $(pwd)/conf:/cuckoo/conf blacktop/cuckoo web`

##### Now Navigate To

-	With [Docker for Mac](https://docs.docker.com/engine/installation/mac/) : `http://localhost`
-	With [docker-machine](https://docs.docker.com/machine/) : `http://$(docker-machine ip)`
-	With [docker-engine](https://docker.github.io/engine/installation/) : `$(docker inspect -f '{{ .NetworkSettings.IPAddress }}' cuckoo-web)`

![cuckoo-submit](https://github.com/blacktop/docker-cuckoo/raw/master/docs/img/2.0/submit.png)  
![cuckoo-config](https://github.com/blacktop/docker-cuckoo/raw/master/docs/img/2.0/config.png)  
![cuckoo-dashboard](https://github.com/blacktop/docker-cuckoo/raw/master/docs/img/2.0/dashboard.png)
