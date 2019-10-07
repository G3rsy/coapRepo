# coapRepo
Tcp/udp (dtls) client/server communication with json and cbor 

This is a little client/server coap implementation with tcp and udp(dtls) standard connection.
Messages can be coded in tcp or cbor object.

Client requires two arguments to work:
  - Connection type: tcp/udp;
  - Endcoding type: json/cbor;
  
Server only have to be builded and runned.
