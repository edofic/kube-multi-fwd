package fwd

import (
	"errors"
	"log"
	"net"
)

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Proxy(stream Proxy_ProxyServer) error {
	var target net.Conn
	rawReq, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("received request", rawReq)
	switch req := rawReq.Req.(type) {
	case *ProxyRequest_Connect:
		log.Println("connecting", req)
		target, err = net.Dial("tcp", req.Connect.Target)
		if err != nil {
			log.Println("error connecting", err)
			return err
		}
		defer target.Close()
		err = stream.Send(&ProxyResponse{Res: &ProxyResponse_Connected{}})
		if err != nil {
			log.Println("error sending connected response", err)
			return err
		}
	default:
		return errors.New("unknown request")
	}
	go func() {
		defer target.Close()
		defer func() {
			stream.Send(&ProxyResponse{Res: &ProxyResponse_Eof{}})
		}()
		for {
			buf := make([]byte, 32*1024)
			n, err := target.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}
			chunk := buf[:n]
			err = stream.Send(
				&ProxyResponse{
					Res: &ProxyResponse_Chunk{
						Chunk: chunk,
					},
				},
			)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	for {
		rawReq, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		switch req := rawReq.Req.(type) {
		case *ProxyRequest_Chunk:
			_, err = target.Write(req.Chunk)
			if err != nil {
				log.Println(err)
				return err
			}
		default:
			log.Println("received request", rawReq)
			return errors.New("unknown request")
		}
	}
}
