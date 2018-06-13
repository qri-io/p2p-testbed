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
RIP(1, 2)
OSPF
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
RMI
Remoting
RPC
HTTP

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



Peer Routing
Swarm
Distributed Record Store
Discovery
Messaging
Naming