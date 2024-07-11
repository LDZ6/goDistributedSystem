package main

import (
	"fmt"
	"goDistributedSystem/protobufTest/prototext"

	"google.golang.org/protobuf/proto"
)

func main() {
	text := &prototext.Test{
		Name:   "panada",
		Weight: []int32{120},
		Height: 180,
		Motto:  "好好学习,天天向上",
	}
	fmt.Println(text)
	data, _ := proto.Marshal(text)
	newText := &prototext.Test{}
	proto.Unmarshal(data, newText)
	fmt.Println(newText)
}
