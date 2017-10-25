package httputils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// **** Request handler ****
type RequestHandler struct {
	getRemoteHostAddress func(r *http.Request) string
	Method               string
	Handler              func(remoteHostAddress string, params *json.RawMessage) (bool, error)
}

// Should be used if there is no reverse proxy i.e. address is not changed
func DefaultHostAddressExtractor(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func sendResponse(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := ServerResponseMessage{Success: success, Message: message}
	json.NewEncoder(w).Encode(response)
}

func (rh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Every path has a mathod
	if r.Method != rh.Method {
		sendResponse(w, false, "Not a valid HTTP method for given path")
		return
	}

	// Lets get all params required for the handler
	var command json.RawMessage
	jsonParser := json.NewDecoder(r.Body)
	jsonParser.Decode(&command)
	remoteHostAddress := rh.getRemoteHostAddress(r)

	// Invoke handler and reply to host based on the response
	success, err := rh.Handler(remoteHostAddress, &command)
	if err != nil {
		sendResponse(w, false, err.Error())
		return
	}

	if success {
		sendResponse(w, true, "Success")
	} else {
		sendResponse(w, false, "Failed executing command")
	}
}

// **** Request multiplexer *****
type RequestMultiplexer struct {
	mux               *http.ServeMux
	hostAddrExtractor func(r *http.Request) string
}

func (rm *RequestMultiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rm.mux.ServeHTTP(w, r)
}

func (rm *RequestMultiplexer) RegisterHandler(path string, method string,
	rh func(string, *json.RawMessage) (bool, error),
) {

	rm.mux.Handle(path, &RequestHandler{
		getRemoteHostAddress: rm.hostAddrExtractor,
		Method:               method,
		Handler:              rh,
	})
}

func (rh *RequestMultiplexer) StartHttpServer(port string) *http.Server {

	addr := fmt.Sprint(":", port)
	h := &http.Server{Addr: addr, Handler: rh}

	go func() {
		if err := h.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	return h
}

func CreateHttpRequestMultiplexer(hostAddressExtractor func(r *http.Request) string) *RequestMultiplexer {
	rh := &RequestMultiplexer{
		mux:               http.NewServeMux(),
		hostAddrExtractor: hostAddressExtractor,
	}

	return rh
}
