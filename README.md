# Visual Lora Gateway connect LoRaWAN™ Network Server 
Visual Machine(Lora Gateway) connect to lora-gateway-bridge for test loar server benchmark or others(NS+NC+AS)

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
	[PROTOCOL](https://github.com/Lora-net/packet_forwarder/blob/master/PROTOCOL.TXT)

### Thanks
	Thanks to [brocaar](https://github.com/brocaar/lora-gateway-bridge)
	