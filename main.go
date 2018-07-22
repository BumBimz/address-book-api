package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const HOST = "localhost:5210"

type AddressBook struct {
	ID      string
	Name    string
	Address string
	Tel     string
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	var result []AddressBook
	session, _ := mgo.Dial(HOST)
	defer session.Close()
	c := session.DB("address").C("books")
	c.Find(bson.M{}).All(&result)
	for _, v := range result {
		fmt.Fprintf(w, "ID: "+v.ID+"\n")
		fmt.Fprintf(w, "Name: "+v.Name+"\n")
		fmt.Fprintf(w, "Address: "+v.Address+"\n")
		fmt.Fprintf(w, "Tel: "+v.Tel+"\n")
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var address AddressBook
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&address)
		if err != nil {
			panic(err)
		}
		fmt.Println(address.Name)
		session, _ := mgo.Dial(HOST)
		defer session.Close()
		c := session.DB("address").C("books")
		c.Insert(&address)
		fmt.Fprintf(w, "ID: "+address.ID+"\n")
		fmt.Fprintf(w, "Name: "+address.Name+"\n")
		fmt.Fprintf(w, "Address: "+address.Address+"\n")
		fmt.Fprintf(w, "Tel: "+address.Tel+"\n")
	} else {
		fmt.Fprintf(w, "Not Found API\n")
	}
}

func manageRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		updateRecord(w, r)
	} else if r.Method == "GET" {
		getDetail(w, r)
	} else if r.Method == "DELETE" {
		deleteRecord(w, r)
	}
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	pathUrl := strings.TrimPrefix(r.URL.Path, "/record/")
	var parameter []string
	parameter = strings.Split(pathUrl, "/")
	if (len(parameter) == 2) && (parameter[1] == "update") {
		id := parameter[0]
		session, _ := mgo.Dial(HOST)
		defer session.Close()
		c := session.DB("address").C("books")
		var address AddressBook
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&address)
		if err != nil {
			panic(err)
		}
		address.ID = id
		err = c.Update(bson.M{"id": id}, bson.M{"$set": &address})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "ID: "+address.ID+"\n")
		fmt.Fprintf(w, "Name: "+address.Name+"\n")
		fmt.Fprintf(w, "Address: "+address.Address+"\n")
		fmt.Fprintf(w, "Tel: "+address.Tel+"\n")
	}
	fmt.Fprintf(w, "OK\n")
}

func getDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/record/")
	if id != "" {
		session, _ := mgo.Dial(HOST)
		defer session.Close()
		c := session.DB("address").C("books")
		result := AddressBook{}
		c.Find(bson.M{"id": id}).One(&result)
		if result != (AddressBook{}) {
			fmt.Fprintf(w, "ID: "+result.ID+"\n")
			fmt.Fprintf(w, "Name: "+result.Name+"\n")
			fmt.Fprintf(w, "Address: "+result.Address+"\n")
			fmt.Fprintf(w, "Tel: "+result.Tel+"\n")
		} else {
			fmt.Fprintf(w, "Not Found DATA\n")
		}
	} else {
		fmt.Fprintf(w, "Not Found API\n")
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/record/")
	if id != "" {
		session, _ := mgo.Dial(HOST)
		defer session.Close()
		c := session.DB("address").C("books")
		result := AddressBook{}
		c.Find(bson.M{"id": id}).One(&result)
		err := c.Remove(bson.M{"id": id})
		if err == nil {
			fmt.Fprintf(w, "Remove\n")
			fmt.Fprintf(w, "ID: "+id+"\n")
			fmt.Fprintf(w, "Name: "+result.Name+"\n")
			fmt.Fprintf(w, "Address: "+result.Address+"\n")
			fmt.Fprintf(w, "Tel: "+result.Tel+"\n")
		} else {
			fmt.Fprintf(w, "Not Found DATA\n")
		}
	} else {
		fmt.Fprintf(w, "Not Found API\n")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.HandleFunc("/record", getRecord)
	http.HandleFunc("/create", create)
	http.HandleFunc("/record/", manageRecord)
	http.ListenAndServe(":3000", nil)
}

//address book api
//id,name,address,tel

// /create => create record => POST
// /record => get all records => GET
// /record/id GET
// /record/id/update PUT
// /record/id/ DEL
