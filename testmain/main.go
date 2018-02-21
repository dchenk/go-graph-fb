package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dchenk/go-graph-fb"
)

func main() {

	if len(os.Args) < 3 {
		log.Println("Missing the function and/or access-token arguments!")
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

	resp, err := fb.Req("GET", "me", accessToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	me := new(fb.GraphResponseMe)
	err = fb.ReadResponse(resp, me)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request went through good!")
	fmt.Printf("%+v \n", me)

}

func mePost(accessToken string) {

	resp, err := fb.Req("POST", "me", accessToken, nil, &fb.ParamStrStr{"method", "GET"})
	if err != nil {
		log.Fatal(err)
	}

	me := new(fb.GraphResponseMe)
	err = fb.ReadResponse(resp, me)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request went through good!")
	fmt.Printf("%+v \n", me)

}
