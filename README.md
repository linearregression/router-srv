# Router Server

Router server provides global routing capabilities. It is not a proxy. It has the selector semantics 
and offloads actual routing to smart clients. Its job is to maintain routing information about the 
environment and route appropriately.

The router also provides weighted and label based routing.

<p align="center">
  <img src="https://github.com/micro/go-platform/blob/master/doc/router.png" />
</p>

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/router-srv
	router-srv --database_url="root:root@tcp(192.168.99.100:3306)/router"
	```

	OR as a docker container

	```shell
	docker run microhq/router-srv --database_url="root:root@tcp(192.168.99.100:3306)/router" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Router server implements the following RPC Methods

Router
- Stats
- Select
- SelectStream

Label
- Read
- Create
- Update
- Delete
- Search
