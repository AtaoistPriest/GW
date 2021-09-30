# GW
a simple VPN demo

This is a simple VPN demo based on UDP. There is no encrypted transmission, no certificate authentication, and no key exchange. It is mainly used to verify the forwarding performance of the machine. There is a log management system and a configuration manager in this demo.Demo supports AF_packet(1,2,3) and pcap packet capture

You should edit cfg first.
Then
```
go build gw.go
sudo ./gw
```
