package main

import (
	"encoding/binary"
	"fmt"
	"github.com/evacchi/proto-example/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if os.Args[1] == "r" {
		read()
	} else {
		write()
	}
}

const numEntries = 1000

func read() {
	buf := make([]byte, 8)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	i := binary.LittleEndian.Uint64(buf)
	buf = make([]byte, i)
	start := time.Now()
	_, err = os.Stdin.Read(buf)
	elapsed := time.Since(start)
	log.Printf("Read: %s elapsed", elapsed)

	if err != nil {
		log.Fatal(err)
	}
	var m pb.AddressBook
	start = time.Now()
	err = proto.Unmarshal(buf, &m)
	elapsed = time.Since(start)
	log.Printf("Unmarshalling: %s elapsed", elapsed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.People[numEntries-1].Name)
}

func write() {
	start := time.Now()
	p := pb.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.Person_HOME},
		},
	}
	book := &pb.AddressBook{}
	for i := 0; i < numEntries; i++ {
		p1 := p
		p1.Name = p1.Name + " " + strconv.Itoa(i)
		book.People = append(book.People, &p1)
	}
	elapsed := time.Since(start)
	log.Printf("Generating: %s elapsed", elapsed)
	start = time.Now()
	bytes, err := proto.Marshal(book)
	if err != nil {
		log.Fatal(err)
	}
	elapsed = time.Since(start)
	log.Printf("Marshalling: %s elapsed", elapsed)

	sz := uint64(len(bytes))
	_, err = os.Stdout.Write(binary.LittleEndian.AppendUint64(nil, sz))

	start = time.Now()
	n, err := os.Stdout.Write(bytes)
	elapsed = time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Written: %d bytes, %s elapsed", n, elapsed)
}
