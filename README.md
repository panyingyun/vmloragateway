### LoRaWAN Packet Simulator (LoRaWAN Node + LoRaWAN Gateway)
  
Virtual LoRaWAN Node and LoRa Gateway connect LoRaWAN™ Network Server 

* 1 Gateway Status Packet (OK)
* 2 Gateway HeartBeat Packet (OK)
* 3 Node Join Packet (ToDo)
* 4 Node Send Packet (ToDo)
* 5 Node Receive Packet (ToDo)


### Usage

1. set gateway id to replace default gateway id

	```shell
	vmloragateway -gw F1E2D3C4B5A69999
	```

2. set more parms to replace default params by config file

	```shell
	vmloragateway -c gateway.conf -gw F1E2D3C4B5A69999
	```
	
### Support features 

* PUSH_DATA/PUSH_ACK 
* PULL_DATA/PULL_ACK
* PULL_RESP/TX_ACK

### PROTOCOL
[Gateway PROTOCOL](https://github.com/Lora-net/packet_forwarder/blob/master/PROTOCOL.TXT)

### Thanks
Thanks to [brocaar](https://github.com/brocaar/lora-gateway-bridge)

	
