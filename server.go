package main

import (
	"sync"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/2tvenom/cbor"
	"github.com/go-ocf/go-coap"
	"github.com/pion/dtls"
	"log"
	"time"
)

type MyOBJ struct {
	IntArr []int
	StrArr []string
}

func (k MyOBJ) logPrintMyOBJ(){
	log.Println("Int array:", k.IntArr, "String array:", k.StrArr)
}

func timestamp(w coap.ResponseWriter, req *coap.Request) {

	resp := w.NewResponse(coap.Content)
	resp.SetOption(coap.ContentFormat, coap.TextPlain)
	resp.SetPayload([]byte(time.Now().String()))

	log.Println("Timestamp transmission: ",  string(resp.Payload()))

	//send timestamp to client
	if err := w.WriteMsg(resp); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func receiveObject(w coap.ResponseWriter, req *coap.Request) {
	msg := req.Msg
	var obj MyOBJ

	//unmarshal and show of message payload
	switch msg.PathString() {
		case "json":
			err := json.Unmarshal(msg.Payload(), &obj)
			if (err != nil) {
				log.Println("Error in Unmarshaling of msg", err)
			} else {
				obj.logPrintMyOBJ()
			}

		case "cbor":
			var buff bytes.Buffer
			encoder := cbor.NewEncoder(&buff)
			ok, err := encoder.Unmarshal(msg.Payload(), &obj)
			if(!ok){
				log.Println("Error in UNnmarshaling of msg", err)
			}else{
				obj.logPrintMyOBJ()
			}
	}
}

func main() {

	//setting for handle the goroutin 
	var wg sync.WaitGroup
	wg.Add(2)

	//setting of various handler
	mux := coap.NewServeMux()
	mux.Handle("timestamp", coap.HandlerFunc(timestamp))
	mux.Handle("json", coap.HandlerFunc(receiveObject))
	mux.Handle("cbor", coap.HandlerFunc(receiveObject))

	log.Println("Starting Server")

	//running of the tcp server and udo(dtls) server
	go func(){
		defer wg.Done()
		log.Fatal(coap.ListenAndServe("tcp", ":5689", mux))
	}()

	go func(){
		defer wg.Done()
		config := dtls.Config{
			
			PSK: func(hint []byte) ([]byte, error) {
				fmt.Printf("Client's hint: %s \n", hint)
				return []byte{0xBC, 0xC1, 0x23}, nil
			},
			PSKIdentityHint: []byte("Pion DTLS Client"),
			CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
		}

		log.Fatal(coap.ListenAndServeDTLS("udp", ":5688", &config, mux))
	}()

	//waiting for termination of goroutin
	wg.Wait()

}