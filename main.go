package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"os/user"
	"path/filepath"
)

// rpcRequest represents a RPC request.
// rpcRequest implements the io.ReadWriteCloser interface.
type rpcRequest struct {
	r    io.Reader     // holds the JSON formated RPC request
	rw   io.ReadWriter // holds the JSON formated RPC response
	done chan bool     // signals then end of the RPC request
}

// NewRPCRequest returns a new rpcRequest.
func NewRPCRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

// Read implements the io.ReadWriteCloser Read method.
func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *rpcRequest) Call() io.Reader {
	go jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}

func main() {

	var conf *Config

	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}

	// default config location
	configFile := filepath.Join(usr.HomeDir, ".factom-anchoring/config.yaml")

	// check if custom config location passed as flag
	flag.StringVar(&configFile, "c", configFile, "config.yaml path")

	flag.Parse()

	log.Printf("Using config: %s\n", configFile)

	// load config
	if conf, err = NewConfig(configFile); err != nil {
		log.Fatal(err)
	}

	log.Printf("Factomd endpoint: %s\n", conf.Factom.Endpoint)

	// UI
	http.Handle("/", http.FileServer(http.Dir("ui/build")))

	// API endpoint
	http.HandleFunc("/v2", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})

	log.Printf("Starting JSON-RPC API at http://localhost:8082\n")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
