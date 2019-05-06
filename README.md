# Netpupper
## What?
<img src="https://i.imgur.com/TXJRKIs.png" width="500">


Netpupper is a prototype for proactive network testing.

## Why?
Traditional network monitoring is only so effective at discovering problems and is fairly poor at identifying transient issues.

This is because, typically, network monitoring products act somewhat "out of band", and operate by periodically
sending very short and low intensity "ping" style tests (often times, literally ICMP.)

These tests are useful for uptime, health, and failure monitoring, but are not useful when troubleshooting more advanced
problems where proof of capability (or culpability...) is the key requirement.

## How?
An alternate method is proposed where nodes can be scattered throughout the network as either clients or servers.

Clients register themselves to a central _controller_ or _scheduler_ server which uses an an API to initiate tests.

The tests are more intense, longer running, and closer to real world traffic patterns.

Test results are stored in a database along with associated _tags_ to produce a high level view of the network health
from the perspective of a real client or application that traverses it.

## Who is this for?
Netpupper was designed primarily for engineers operating large networks with a lot of variables to their design and 
implementation, where "as-built" benchmarks are not sufficient to prove innocence when problems arise.

A common use case would be a WAN backbone utilizing various carriers, transport mechanisms, and network hardware -
netpupper would be useful to ensure the backbone was operating within required parameters with a very high degree of
assurance. 

## Prototype?
This code is firmly in the __prototype__ phase. While it works, it is missing many features and is no where near
stable enough to run in production unless you are willing to squash your own bugs.
It also features NO SECURITY FEATURES!

It was written primarily to determine:

 - The effectiveness of Golang as a language for developing network-related testing
 - Whether a client/server architecture makes sense for an application of this type
 
The conclusion was that the prototype was successful and may be, one day, turned into a stable build.


## What's with the name?
A pupper is a dog, so it's like the Network Dog, sniffing out problems. It sounds lame but c'mon, have you seen how many
good names are taken on github already?

## Usage
### Build
```
go get github.com/adambaumeister/netpupper.git
cd github.com/adambaumeister/
go build -o netp
```

### Run
See the daemon and daemon_server files for YAML examples.
```bash
# client
./netp -daemon -config daemon.yml

# Server
./netp -daemon -config daemon_server.yml
```


