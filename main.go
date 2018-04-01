package main
import (
	"github.com/stianeikeland/go-rpio"
	"github.com/gorilla/mux"
	"time"
	"io"
	"net/http"
	"github.com/GeertJohan/go.rice"
)

// create a map/dictionary to map the door # to the pin #
// This could also be a name or any string
var pinmap map[string]int

func main() {
    pinmap = make(map[string]int) // have to "make" the map in memory
    pinmap["1"] = 2 // set door 1 to pin 2

    //Create a new Gorilla Router
    r := mux.NewRouter()

    //Add the route to push the button.
    //Notice the {id} part. This will be parsed into id returned from mux.Vars
    r.HandleFunc("/push/{id}",PushButton).Methods("POST")

    //Create a "box" to "rice the static content
    r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("static").HTTPBox()))


    http.ListenAndServe(":8002",r) // start the web server on port 8001
}

func PushButton(w http.ResponseWriter, r *http.Request) {
    //Grab the variables
    vars := mux.Vars(r)
    if p, ok := pinmap[vars["id"]]; ok {
        // All our same code here, we are just using a variable now
    	err := rpio.Open()
    	if err != nil {
       	    panic(err)
        }
        pin := rpio.Pin(p)
        pin.Output()
        pin.Low()
        time.Sleep(time.Millisecond * 100)
        pin.High()
        time.Sleep(time.Millisecond * 100)
        rpio.Close()
        io.WriteString(w,"pushed!")
    } else {
        // Pin was not found in pinmap. Set status code 404
        w.WriteHeader(http.StatusNotFound)
    	io.WriteString(w,"unknown pin!")
    }

}
