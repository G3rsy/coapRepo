package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/2tvenom/cbor"
	"github.com/go-ocf/go-coap"
	"github.com/pion/dtls"
	"log"
	"os"
)


var msgPayload []byte
var co *coap.ClientConn
type myOBJ struct {
	IntArr []int
	StrArr []string
}

func (k myOBJ) logPrintMyOBJ(){
	log.Println("Int array:", k.IntArr, "String array:", k.StrArr)
}

//generation json and cbor object
func jsonMessage() []byte{
	something := myOBJ{
		IntArr: []int{1,2,3},
		StrArr:   []string{"string1", "string2", "string3"}}

	someB,err := json.Marshal(something)
	if(err!=nil){
		log.Fatal("Json conversion failed", err)
	}
	return someB
}
func cborMessage() []byte{
	something := myOBJ{
		IntArr: []int{1,2,3},
		StrArr:   []string{"string1", "string2", "string3"}}

	var buff bytes.Buffer
	encoder := cbor.NewEncoder(&buff)
	ok, err := encoder.Marshal(something)
	if(!ok){
		log.Fatal("Error in cbor translation", err)
	}

	return buff.Bytes()
}

//send request to server and print response
func clientTimestamp(conn *coap.ClientConn){
	path := "timestamp"
	
	resp, err := conn.Get(path)
	
	if err != nil {
		log.Println("Error sending request: ", err)
	}

	log.Printf("Received timestamp from server: %s",	resp.Payload())
}

//send object to server
func sendObject(conn *coap.ClientConn){

	//message generation
	toSend := conn.NewMessage(coap.MessageParams{Payload:msgPayload})
	toSend.SetPathString(os.Args[2])

	log.Println("Transmitting to Server", os.Args[1], "object")

	if err := conn.WriteMsg(toSend); err != nil {
		log.Println("Cannot send message: ", err)
	}
}

func main() {
	
	if(len(os.Args) != 3){
		log.Fatal("Error, program need 2 argoument: tcp/udp json/cbor")
	}
	

	var err error


	//two port for two different type of connection
	switch os.Args[1]{
	case "tcp":
		co, err = coap.Dial("tcp", "localhost:5689")
		break;

	case "udp":
		config := dtls.Config{
			PSK: func(hint []byte) ([]byte, error) {
				fmt.Printf("Client's hint: %s \n", hint)
				return []byte{0xBC, 0xC1, 0x23}, nil
			},
			PSKIdentityHint: []byte("Pion DTLS Client"),
			CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
		}

		//dtls connection starting
		co, err = coap.DialDTLS("udp", "localhost:5688", &config)

		break;
	}
	
	
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	//message payload generation
	switch os.Args[2]{
	case "json":
		msgPayload = jsonMessage()
	case "cbor":
		msgPayload = cborMessage()
	}


	//test
	var i int
	for i=0; i<2000; i++{
		if(i%3 == 0){
			clientTimestamp(co)
		}else{
			sendObject(co)
		}
	}

	log.Println("Client has finished")
}

