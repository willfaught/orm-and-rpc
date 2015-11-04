package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/facebookgo/grace/gracenet"
	"github.com/kr/pretty"
	tigertonic "github.com/rcrowley/go-tigertonic"
	"github.com/company/gossie/src/mockgossie"
	"github.com/willfaught/company/prototype/company/offer"
	"github.com/willfaught/company/prototype/company/offers"
)

// ExampleBasic shows basic use.
func ExampleBasic() {
	log.Println("\n***** ExampleBasic *****\n")

	// Server
	var server = offers.MustNewServer(configure("Server"), offers.Offers{Configuration: configure("Server")}, ":4010")
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	// Client
	var client = offers.MustNewClient(configure("Client"), ":4010").ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)}).(offers.Client)

	// Use
	m, err := client.New(offer.Offer{Name: "nerddomo"})
	if err != nil {
		panic(err)
	}
	log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "test"
	m, err = client.Set(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client.Delete(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// Client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Server
	if err := server.Stop(); err != nil {
		panic(err)
	}
}

// ExampleGracefulRestart shows enabling graceful restarts.
func ExampleGracefulRestart() {
	log.Println("\n***** ExampleGracefulRestart *****\n")

	// Server
	var gn = &gracenet.Net{}
	var listener, err = gn.Listen("tcp", ":4545")
	if err != nil {
		panic(err)
	}
	var server = offers.Server{Listener: listener, Server: rpc.NewServer()}
	server.RegisterName("Marketing", offers.Receiver{Configuration: configure("Server"), Interface: offers.Offers{Configuration: configure("Server")}})
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()
	var done = make(chan struct{})
	go func() {
		var c = make(chan os.Signal, 10)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGUSR2)
		for {
			var s = <-c
			switch s {
			case syscall.SIGTERM:
				log.Println("Received SIGTERM\n")
				signal.Stop(c)
				done <- struct{}{}
				return
			case syscall.SIGUSR2:
				log.Println("Received SIGUSR2\n")
				if _, err := gn.StartProcess(); err != nil {
					log.Fatalln("Error:", err)
				}
			}
		}
	}()

	// Client
	var client = offers.MustNewClient(configure("Client"), ":4545").ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)}).(offers.Client)

	// Use
Main:
	for {
		var m = offer.Offer{Interface: client, Name: "nerddomo"}
		m, err := m.New()
		if err != nil {
			panic(err)
		}
		log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
		m, err = client.Get(m.ID)
		if err != nil {
			panic(err)
		}
		log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
		m.Name = "bar"
		m, err = m.Set()
		if err != nil {
			panic(err)
		}
		log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
		m, err = m.Delete()
		if err != nil {
			panic(err)
		}
		log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))
		log.Println()
		select {
		case <-done:
			break Main
		case <-time.After(time.Second * 10):
		}
	}

	// Client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Server
	if err := server.Stop(); err != nil {
		panic(err)
	}
}

// ExampleHTTPJSON shows using JSON over HTTP.
func ExampleHTTPJSON() {
	log.Println("\n***** ExampleHTTPJSON *****\n")

	// Listener
	var listener, err = net.Listen("tcp", ":5656")
	if err != nil {
		panic(err)
	}

	// RPC server
	var rserver = rpc.NewServer()
	rserver.RegisterName("Marketing", offers.Receiver{Configuration: configure("Server"), Interface: offers.Offers{Configuration: configure("Server")}})

	// Gob server
	var mserver = offers.Server{Listener: listener, Server: rserver}
	go func() {
		if err := mserver.Start(); err != nil {
			panic(err)
		}
	}()

	// JSON server
	var mux = http.NewServeMux()
	mux.Handle("/", offers.NewJSONRPCHandler(rserver))
	var hserver = &http.Server{Addr: ":7777", Handler: mux}
	go func() {
		if err := hserver.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// RPC client
	rclient, err := jsonrpc.Dial("tcp", ":5656")
	if err != nil {
		panic(err)
	}
	var client = offers.Client{Client: rclient, Configuration: configure("Client"), Name: "Marketing"}.ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)}).(offers.Client)

	// Use RPC client
	m, err := client.New(offer.Offer{Name: "nerddomo"})
	if err != nil {
		panic(err)
	}
	log.Printf("RPC New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("RPC Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "test"
	m, err = client.Set(m)
	if err != nil {
		panic(err)
	}
	log.Printf("RPC Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client.Delete(m)
	if err != nil {
		panic(err)
	}
	log.Printf("RPC Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// HTTP client
	var call = func(method string, args map[string]interface{}) interface{} {
		b, err := json.Marshal(map[string]interface{}{
			"id":     client.Context.ID,
			"method": "Marketing." + method,
			"params": []interface{}{
				map[string]interface{}{
					"Arguments": args,
					"Context":   client.Context,
				},
			},
		})
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:7777", strings.NewReader(string(b)))
		if err != nil {
			panic(err)
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		var body struct {
			Error  interface{} `json:"error"`
			ID     interface{} `json:"id"`
			Result interface{} `json:"result"`
		}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			panic(err)
		}
		if body.Error != nil {
			panic(err)
		}
		return body.Result
	}

	// Use HTTP client
	log.Printf("HTTP New: Got: %# v\n\n", call("New", map[string]interface{}{"M": offer.Offer{Name: "nerddomo"}}))

	// RPC client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Marketing server
	if err := mserver.Stop(); err != nil {
		panic(err)
	}

	// Listener
	if err := listener.Close(); err != nil {
		panic(err)
	}
}

// ExampleMethod shows using a method.
func ExampleMethod() {
	log.Println("\n***** ExampleMethod *****\n")

	// Server
	var server = offers.MustNewServer(configure("Server"), offers.Offers{Configuration: configure("Server")}, ":4020")
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	// Client
	var client = offers.MustNewClient(configure("Client"), ":4020").ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)}).(offers.Client)

	// Use
	var m = offer.Offer{Interface: client, Name: "nerddomo"}
	m, err := m.New()
	if err != nil {
		panic(err)
	}
	log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "bar"
	m, err = m.Set()
	if err != nil {
		panic(err)
	}
	log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = m.Delete()
	if err != nil {
		panic(err)
	}
	log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// Client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Server
	if err := server.Stop(); err != nil {
		panic(err)
	}
}

// ExampleProxy shows a chain of servers and clients.
func ExampleProxy() {
	log.Println("\n***** ExampleProxy *****\n")

	// Real
	var realServer = offers.MustNewServer(configure("Real Server"), offers.Offers{Configuration: configure("Server")}, ":4030") // Use s as the service.
	go func() {
		if err := realServer.Start(); err != nil {
			panic(err)
		}
	}()
	var realClient = offers.MustNewClient(configure("Real Client"), ":4030")

	// Proxy
	var proxyServer = offers.MustNewServer(configure("Proxy Server"), realClient, ":4031") // Use realClient as the service.
	go func() {
		if err := proxyServer.Start(); err != nil {
			panic(err)
		}
	}()
	var proxyClient = offers.MustNewClient(configure("Proxy Client"), ":4031").ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)}).(offers.Client)

	// Use
	m, err := proxyClient.New(offer.Offer{Name: "nerddomo"})
	if err != nil {
		panic(err)
	}
	log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = proxyClient.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "test"
	m, err = proxyClient.Set(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = proxyClient.Delete(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// Proxy
	if err := proxyClient.Close(); err != nil {
		panic(err)
	}
	if err := proxyServer.Stop(); err != nil {
		panic(err)
	}

	// Real
	if err := realClient.Close(); err != nil {
		panic(err)
	}
	if err := realServer.Stop(); err != nil {
		panic(err)
	}
}

// ExampleVersion shows using incompatible versions of Marketinges in tandem.
func ExampleVersion() {
	log.Println("\n***** ExampleVersion *****\n")

	// Server
	var listener, err = net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	var server = offers.Server{Listener: listener, Server: rpc.NewServer()}
	server.RegisterName("Marketing", offers.Receiver{Configuration: configure("Server 1"), Interface: offers.Offers{Configuration: configure("Server 1")}})
	server.RegisterName("Marketing_v2", offers.Receiver{Configuration: configure("Server 2"), Interface: offers.Offers{Configuration: configure("Server 2")}} /* v2.Receiver{Configuration: [...], Service: [...]} */)
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	// Clients
	client, err := rpc.Dial("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	var client1 = offers.Client{Client: client, Configuration: configure("Client 1"), Name: "Marketing"}.ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)})
	var client2 = offers.Client{Client: client, Configuration: configure("Client 2"), Name: "Marketing_v2"}.ForContext(offers.Context{ID: tigertonic.RandomBase62String(8)})

	// Use 1
	m, err := client1.New(offer.Offer{Name: "nerddomo"})
	if err != nil {
		panic(err)
	}
	log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client1.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "test"
	m, err = client1.Set(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client1.Delete(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// Use 2
	m, err = client2.New(offer.Offer{Name: "nerddomo"})
	if err != nil {
		panic(err)
	}
	log.Printf("New: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client2.Get(m.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("Get: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m.Name = "test"
	m, err = client2.Set(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Set: Got: %# v\n\n", pretty.Formatter(simplify(m)))
	m, err = client2.Delete(m)
	if err != nil {
		panic(err)
	}
	log.Printf("Delete: Got: %# v\n\n", pretty.Formatter(simplify(m)))

	// Client
	if err := client.Close(); err != nil {
		panic(err)
	}

	// Server
	if err := server.Stop(); err != nil {
		panic(err)
	}
}

func configure(name string) offers.Configuration {
	return offers.Configuration{Name: name, OfferRepository: offer.NewRepository(mockgossie.NewMockConnectionPool())}
}

func main() {
	log.SetFlags(0)
	if false {
		ExampleGracefulRestart()
	} else {
		ExampleBasic()
		//ExampleHTTPJSON()
		//ExampleMethod()
		//ExampleProxy()
		//ExampleVersion()
	}
}

func simplify(o offer.Offer) offer.Offer {
	o.Interface = nil
	return o
}
