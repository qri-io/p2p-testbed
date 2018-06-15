# p2p testbed

p2p testbed uses request tracing to demonstrate distributed techniques and algorithms using [libp2p](https://libp2p.io) and [request tracing](https://opentracing.io)

# Kasey's notes while tryin to learn about libp2p:

Okay folks, Kasey (aka @ramfox) here. I'm trying to learn how the [libp2p](https://libp2p.io/) package works, and how p2p works in general. 
For those who don't know, Qri is built on top of the distributed web. That's right muggles, you didn't even realize it, but just by downloading the [Qri electron app](../frontend) and opening the app or running `qri connect` from the terminal, you are joining the rest of us wizards on the distributed web. Not only is Qri on the distributed web, but Qri specifically uses IPFS, aka the InterPlanitary File System. They have a specification and implimentation of libp2p, and imma learn how it works and what it does.

I started off checking out the [libp2p specification](http://www.github.com/libp2p):
You should go take a read, this README is more or less a response to it:

It begins with a _why do we need libp2p_: they need a way for peers to communicate with each other over IPFS, but they also can't make assumtions about which networks/protocols those peers will use. Basically, let the devs choose how they want to party, and build a tool that is flexible enough to let different apps party the way they want: "In essence, a peer using libp2p should be able to communicate with another peer using a variety of different transports, including connection relay, and talk over different protocols, negotiated on demand." [source](https://github.com/libp2p/specs/blob/master/1-introduction.md#11-motivation)

K, made it through section 1.1, one to 1.2... And that's where I got lost. Cause guess what! I don't remember anything from the networking class I took in college half a decade ago! Section 1.2 started talking about the goals for libp2p, and suddenly I was lost. They started listing all these different network protocols, like TCP, UDP, SCTP, RNG, TSM, SKT, and if you are nodding in understanding while reading those then you are just as lost as I am cause those last three aren't network protocols they are League of Legends teams.

So girlfriend's gonna go through all these network buzz words they are using so she can actually learn what's up. And if you are having trouble diving into this topic, maybe it will help you too :)

(Don't want to bog down my time rn with researching all of these right away, I'm pretty sure that most of it will come up organically as I read into the spec and implimentation)

### 2.2.1 Establishing the physical link
Ethernet
Wi-Fi
Bluetooth
USB
### 2.2.2 Addressing a machine or process
IPv4 - IP version 4 -> 32-bit address, gives us 4.2 billion ip addresses, made at a time when the internet was just a thing for academics
IPv6 - IP version 6 -> 128-bit address, gives us 3.4×1038 addresses ip addresses
Hidden addressing, like SDP
### 2.2.3 Discovering other peers or services
ARP - Address Resolution Protocol is a protocol for mapping an IP address to a physical machine address that is recognized in the local network ([Learn Address Resolution Protocol](https://www.youtube.com/watch?v=ULpPIVln6nI)) -> physical address known as a MAC address (Media Access Control) ([What is a MAC address](https://www.youtube.com/watch?v=UrG7RTWIJak))
DHCP -> Dynamic Host Configuration Protocol -> the protocol that allows you to request an ip address from a DHCP server, eg your computer requesting an ip address from your router, or your modem requesting an ip address from your ISP ([What is DHCP?](https://www.lifewire.com/what-is-dhcp-2625848))
DNS - Domain Name System -> [Introduction to the Domain Name System](https://www.lifewire.com/introduction-to-domain-name-system-817512)
Onion -> encrypted anonymous routing, used by the TOR browser (The Onion Routing) -> [Onion Routing - Computerphile](https://www.youtube.com/watch?v=QRYzre4bf7I)
### 2.2.4 Routing messages through the network
RIP(1, 2) - older protocol -  Routing Information Protocol used to route data packets by finding the best hop count
OSPF - Open Shortest Path First - routing protocol for IP networks
BGP
PPP
Tor
I2P
cjdns
### 2.2.5 Transport
TCP - Transmission Control Protocol - 'reliable and connection-based' -> needs to set up a connection a 'three-way handshake'
UDP - User Datagram Protocol - 'lightweight and connectionless' -> unreliable, no error recovery or lost packet recovery, 
[UDP and TCP: Comparison of Transport Protocols](https://www.youtube.com/watch?v=Vdc8TCESIg8)
UDT -> UDP-Based Data Transfer Protocol - "is a high-performance data transfer protocol designed for transferring large volumetric datasets over high-speed wide area networks." [wiki](https://en.wikipedia.org/wiki/UDP-based_Data_Transfer_Protocol)
QUIC -> Quick UDP Internet Connection - [website](https://www.chromium.org/quic)
WebRTC data channel -> Web Real-Time Communication -> [wiki](https://en.wikipedia.org/wiki/WebRTC)
### 2.2.6 Agreed semantics for applications to talk to each other
RMI - Javascript equivalent for remove procedure calls
Remoting
RPC - Remote procedure calls -> way for 
HTTP - > 

What is NAT Traversal:
NAT (Network address translator) traversal is the work around used to compensate the fact that there are more devices trying to use the internet then there are ip addresses for them in IPv4. Watch this video for a good explaination: [How Network Address Translation Works](https://www.youtube.com/watch?v=QBqPzHEDzvo), and may I recommend watching it at 1.25 speed. There are varing protocols for nat traversal

ICE -> Interactive Connectivity Establishment

### WTF is multiplexing (or muxing), and WhyTF do we care?
According to the [Macao Communications Museum website](http://macao.communications.museum/eng/exhibition/secondfloor/moreinfo/2_8_6_Multiplexing.html), which is randomly in english, although I guess that's not too weird: "Multiplexing is the process in which multiple Data Streams, coming from different Sources, are combined and Transmitted over a Single Data Channel or Data Stream."

So basically, a bunch of signals are trying to jam over one 'wire', multiplexing combines these multiple signals into one signal. This concept is not new, It's been around since the telegraph. 

We use it because it reduces the resources needed to transmit multiple signals

In the TCP/IP and OSI models, multiplexing happens at the transport layer, and it lets multiple applications use one network connection simultaneously 

Resourses to learn about networking, super usefulllll:
video - [Networking Crash Course](https://www.youtube.com/watch?v=mgEMGoFIots)
article - [The Internet, Networks, and TCP/IP](http://www.cellbiol.com/bioinformatics_web_development/chapter-1-internet-networks-and-tcp-ip/)



### 4.1 Peer Routing -> how we decide which peers to use for routing messages
4.1.1 kad-routing -> [P2P Networks - Basic Algorithms](https://www.youtube.com/watch?v=kXyVqk3EbwE) (watch at 1.5 speed)

Swarm
Distributed Record Store
Discovery
Messaging
Naming


What is a Distributed Hash Table?
"A distributed hash table is a class of decentralized distributed system that provides a lookup service similar to a hash table: (key, value) pairs are storied in a DHT and any participating node can efficiently retrieve the value associated with a given key."
DHTs characteristically emphasize the following properites:
* Autonomy and decentralization: the nodes collectively form the system without any central coordination
* Fault tolerance: the system should be reliable (in some sense) even with nodes continuously joining, leaving, and failing
* Scalability: the system should function efficiently even with thousands or millions of nodes

A key techinique used to achieve these goals is that any one node needs to coordinate with only a few other nodes in the system - most commonly O(log n) of the n participants

Overlay network: The set of links each node maintains (it's neighbors or routing table), together form the overlay network. The node picks its neighbors according to a certain structure, called the network's topology

ALL DHT TOPOLOGIES SHARE SOME VARIENT OF THE MOST ESSENTIAL PROPERTY: for any key _k_, each node either has a node ID that owns _k_ or has a link to a node whose node ID is _closer_ to _k_. We use a greedy algorithm to forward messages: at each step, forward the message to the neighbor whose ID is closest to k, when there is no such neighbor, we must have arrived at k. Sometimes called key-based routing.

Swarm is an sbstraction layer over connection management
Had no f-ing idea what 'swarm' meant, what it was supposed to do. But then I read this issue: [New and more accurate name for libp2p-swarm](https://github.com/libp2p/js-libp2p-switch/issues/40), and it's amazing how much less confused I am now looking at the swarm docs

The Peer Routing is the mechanism/logic for how we decide what peers to contact, in order to get to the destination peer.
The Distributed Record Store, is the place where the list of peers that we can contact is kept (info includes their peerIDs and IP addressess)
The Swarm (think 'switch' instead) is what opens/closes/manages connections and handles muxing. Once you have a swarm set up, you can open and close streams to other peers.

Yo, reading the docs on [Circuit Relay](https://github.com/libp2p/specs/tree/master/relay), they are so good. Examples of multiaddrs and why they are useful

Distributed Record Store: [Interplanetary Record Store spec](https://github.com/libp2p/specs/blob/master/IPRS.md)

Okay, so after all that context, I'm finally going to look at the p2p-testbed and see if I'm less confused than before.


### testbed.go
Looking at the main func:
what is startAppDash(), I'm guessing it... starts whatever AppDash is.

  startAppDash:
  1) get a new MemoryStore
  2) ListenTCP - listen on any available TCP port locally
  3) Get's the port from the addr -> collectorPort
  4) collectorAdd = ":[collectorPort]" ???
  5) appdash.Newserver(listener, collector) returns a collectorserver
  6) in a go routine, start the server
  7) prints the port/url the webUI will be running
  8) parse url into type URL
  9) start the webUI in a separate goroutine - traceapp.New returns a new traceapp, with a router with a bunch of handlers, listen and serve in a function running in a separate go routine
  11) appdashot.Tracer reports spans to the collector. 

we get a context.Background(), which means we get a context that does not have a endtime
Okay what does Setup(context, peerCount) do? -> returns a slice of references to Peers
  Setup():
  1) creates a span (opentracing.Span)
  2) StartSpanFromContext()? => calls startSpanFromContextWithTracer, which sees if there is already a span associated with Context (parentSpan), if so, it adds it to the list of StartSpanOptions (opts), returns a span
  3) SetTag "peerCount" # of peers
  4) Create a new network -> 
    1) new span, StartSpanFromContext(ctx, "NewNetwork")
    2) make slice of peers, underlying array # of peerCount
      For each num, create RandPeerNetParams() (from the libp2p/go-testutil package)
      PeerNetParams include an ID, PrivKey, PubKey, and Addr (multiaddr)
      Makes new PeerNetParams, Addr is the ZeroLocalTCPAddress, random Priv and Pub keys are generated, ID is generated from PubKey, and PeerNetParams are returned
    3) Create a new peer from the ctx and params
    Peers have: Id, Host, Peerstore, ctx, and some state that holds info (map[string]interface{})
    Peerstore => https://github.com/libp2p/go-libp2p-peerstore/blob/master/peerstore.go
    makes a new peerstore, adds private and public key to peerstore
    makeBasicHost(context, ID, slice of multiaddrs, with addr to peer first, peerstore)

Okay time ot find out about opentracing, yay learning new things




