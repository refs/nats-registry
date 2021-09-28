## Run this garbage

0. start a nats server: `docker run -p 4222:4222 -ti nats:latest`
1. run this `go run main.go`
2. get the port echoer is running on
3. get registered nodes at the registry: `curl 0.0.0.0:56974/list`
4. register a new node: `curl -X POST 0.0.0.0:56974/register -d 'name=echoer' -d 'address=0.0.0.0:1234'`
5. request the new nodes list and see it updated

Long story short, the registry is subscribed to the subject `register_service` and whenever a message comes it updates its state. The current persistent layer is in memory so it is short-lived. This also forces to use a single registry service, as putting a load balancer on top would result in different data since they are not yet synchronized, and there is no gossip at all.