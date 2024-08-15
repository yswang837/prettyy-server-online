package register

import (
	"log"
	"net"
	"testing"
)

func TestRegisterConsul(t *testing.T) {
	r := NewRegisterConsul()
	if err := r.Init(); err != nil {
		log.Fatalf("init consul fail:%v", err.Error())
	}
	if err := r.Register("wysConsul", 4567); err != nil {
		log.Fatalf("register consul fail:%v", err.Error())
	}
	l, err := net.Listen("tcp", ":4567")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			log.Printf("Ip: %s connected", conn.RemoteAddr().String())
		}()
	}
}
