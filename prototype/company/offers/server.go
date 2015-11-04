// Generated

package offers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	"github.com/willfaught/company/prototype/company/offer"
)

var _ io.ReadWriteCloser = &readWriteCloser{}

type readWriteCloser struct {
	r io.Reader
	w io.Writer
}

func (r *readWriteCloser) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *readWriteCloser) Write(p []byte) (n int, err error) {
	return r.w.Write(p)
}

func (r *readWriteCloser) Close() error {
	return nil
}

// NewJSONRPCHandler makes a JSON-RPC handler for s.
func NewJSONRPCHandler(s *rpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		var b bytes.Buffer
		var codec = jsonrpc.NewServerCodec(&readWriteCloser{r: r.Body, w: &b})
		if err := s.ServeRequest(codec); err != nil {
			panic(err)
		}
		codec.Close()
		io.Copy(w, &b)
	})
}

// Receiver is an RPC receiver.
type Receiver struct {
	// Configuration is the configuration.
	Configuration Configuration
	// Interface is the service.
	Interface Interface
}

// Delete deletes M.
func (r Receiver) Delete(argument struct {
	Arguments struct { // Method arguments
		M offer.Offer // m
	}
	Context Context // Call context
}, result *struct { // Method result
	Offer offer.Offer // Only one non-error result
}) error {
	var err error
	func() { // Panic recovery
		defer func() { // Panic handler
			var value = recover()
			if value == nil {
				return // No panic
			}
			r.Configuration.OnPanic("Delete", argument.Context, value) // Called because Configuration implements it
			err = errors.New(fmt.Sprint(value))                        // Put the panic error into the RPC error result
		}()
		r.Configuration.OnCall("Delete", argument.Context, argument.Arguments)                    // Called because Configuration implements it
		result.Offer, err = r.Interface.ForContext(argument.Context).Delete(argument.Arguments.M) // ForContext is called because Service implements it
		r.Configuration.OnReturn("Delete", argument.Context, *result, err)                        // Called because Configuration implements it
		result.Offer.Interface = nil                                                              // Avoid encoding the Service or Client
	}()
	return err
}

// Get gets the Offer for id.
func (r Receiver) Get(argument struct {
	Context   Context
	Arguments struct {
		Id string
	}
}, result *struct {
	Offer offer.Offer
}) error {
	var err error
	func() {
		defer func() {
			var value = recover()
			if value == nil {
				return
			}
			r.Configuration.OnPanic("Get", argument.Context, value)
			err = errors.New(fmt.Sprint(value))
		}()
		r.Configuration.OnCall("Get", argument.Context, argument.Arguments)
		result.Offer, err = r.Interface.ForContext(argument.Context).Get(argument.Arguments.Id)
		r.Configuration.OnReturn("Get", argument.Context, *result, err)
		result.Offer.Interface = nil
	}()
	return err
}

// New creates M.
func (r Receiver) New(argument struct {
	Context   Context
	Arguments struct {
		M offer.Offer
	}
}, result *struct {
	Offer offer.Offer
}) error {
	var err error
	func() {
		defer func() {
			var value = recover()
			if value == nil {
				return
			}
			r.Configuration.OnPanic("New", argument.Context, value)
			err = errors.New(fmt.Sprint(value))
		}()
		r.Configuration.OnCall("New", argument.Context, argument.Arguments)
		result.Offer, err = r.Interface.ForContext(argument.Context).New(argument.Arguments.M)
		r.Configuration.OnReturn("New", argument.Context, *result, err)
		result.Offer.Interface = nil
	}()
	return err
}

// Set updates M.
func (r Receiver) Set(argument struct {
	Context   Context
	Arguments struct {
		M offer.Offer
	}
}, result *struct {
	Offer offer.Offer
}) error {
	var err error
	func() {
		defer func() {
			var value = recover()
			if value == nil {
				return
			}
			r.Configuration.OnPanic("Set", argument.Context, value)
			err = errors.New(fmt.Sprint(value))
		}()
		r.Configuration.OnCall("Set", argument.Context, argument.Arguments)
		result.Offer, err = r.Interface.ForContext(argument.Context).Set(argument.Arguments.M)
		r.Configuration.OnReturn("Set", argument.Context, *result, err)
		result.Offer.Interface = nil
	}()
	return err
}

// Server is the RPC server.
type Server struct {
	// Server is the wrapped RPC server.
	*rpc.Server
	// Listener is the listener.
	Listener net.Listener
}

// MustNewServer calls NewServer and panics if there is an error.
func MustNewServer(c Configuration, i Interface, address string) Server {
	var server, err = NewServer(c, i, address)
	if err != nil {
		panic(err)
	}
	return server
}

// NewServer makes a Server that uses h and s and listens to address.
func NewServer(c Configuration, i Interface, address string) (Server, error) {
	var listener, err = net.Listen("tcp", address)
	if err != nil {
		return Server{}, err
	}
	var server = rpc.NewServer()
	server.RegisterName("Marketing", Receiver{Configuration: c, Interface: i}) // Default receiver name is "Marketing".
	return Server{Listener: listener, Server: server}, nil
}

// Start starts the server.
func (s Server) Start() error {
	for {
		var conn, err = s.Listener.Accept()
		if err != nil {
			if strings.HasSuffix(err.Error(), "use of closed network connection") { // Standard library badly designed.
				return nil
			}
			return err
		}
		go s.Server.ServeConn(conn)
	}
}

// Stop stops the server.
func (s Server) Stop() error {
	return s.Listener.Close()
}
