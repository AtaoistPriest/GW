# GW
a simple VPN demo

This is a simple UDP based VPN demo. There is no encrypted transmission, no certificate authentication, and no key exchange. It is mainly used to verify the forwarding performance of the machine. In this demo, the log management system and configuration manager support AF by capturing packets_ Pack (1,2,3) and pcap.

You should edit cfg first.
```
go build gw.go
sudo ./gw
```
