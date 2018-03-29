package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dchenk/go-graph-fb"
)

func main() {

	if len(os.Args) < 3 {
		log.Println("Missing the function and/or access-token arguments: give function key first, then token.")
		return
	}

	f := os.Args[1]

	if _, ok := argFuncs[f]; !ok {
		log.Printf("The provided argument %q does not map to a function.\n", f)
		return
	}

	argFuncs[f](os.Args[2])

}

var argFuncs = map[string]func(string){
	"me":      me,
	"me-post": mePost,
}

func me(accessToken string) {

	resp, err := fb.ReqDo("GET", "me", accessToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	me := new(fb.GraphResponseMe)
	if err := fb.ReadResponse(resp, me); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request went through good!")
	fmt.Printf("%+v \n", *me)

}

func mePost(accessToken string) {

	resp, err := fb.ReqDo("POST", "me", accessToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	me := new(fb.GraphResponseMe)
	err = fb.ReadResponse(resp, me)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request went through good!")
	fmt.Printf("%+v \n", *me)

}
