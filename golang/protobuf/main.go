package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
	"sgorecki.me/protobuf/github.com/protocolbuffers/protobuf/sgorecki/me"
)

func main() {
	p := me.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*me.Person_PhoneNumber{
			{Number: "555-4321", Type: me.PhoneType_PHONE_TYPE_HOME},
		},
	}
	out, err := proto.Marshal(&p)
	if err != nil {
		log.Fatalln("Failed to encode person:", err)
	}
	fmt.Printf("%b\n", out)

	p2 := &me.Person{}
	if err := proto.Unmarshal(out, p2); err != nil {
		log.Fatalln("Failed to decode person:", err)
	}
	fmt.Printf("%s\n", p2)
}
