docker build -t coapserver .
docker run -p 5688:5688/udp -p 5689:5689 coapserver
