package main

import (
	"testing"
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"net/url"
	"io"
)

const baseUrl = "http://localhost:3000"
var userUrl = fmt.Sprintf("%s/users", baseUrl)
var itemUrl = fmt.Sprintf("%s/items", baseUrl)
var deleteUrl = fmt.Sprintf("%s/all", baseUrl)


type UserModel struct {
	Name  		string
	Job   		string
	Level 		int
}


type ItemModel struct {
	Name 			string
	Level 			int
	Strength 		int
	Dexterity 		int
	Intelligence 	int
	Vitality 		int
}


func TestDeleteAll(t *testing.T) {
	ClearAll()

	bodyString := ExecuteRequest(http.NewRequest(http.MethodGet, userUrl, nil))
	expected := "[]"
	Compare(t, bodyString, expected)

	bodyString = ExecuteRequest(http.NewRequest(http.MethodGet, itemUrl, nil))
	expected = "[]"
	Compare(t, bodyString, expected)
}


// level(>= 1, <= 20), job(Barbarian, Mage, Hunter)
func TestCreateValidUserModels(t *testing.T) {
	ClearAll()
	defer ClearAll()

	cases := []UserModel {
		UserModel{"Timothy", "Mage", 1},
		UserModel{"Shawn", "Barbarian", 10},
		UserModel{"Anthony", "Hunter", 20},
	}

	for _, c := range cases {
		body := CreateUserModelBody(c)
		ExecuteRequest(http.NewRequest(http.MethodPost, userUrl, body))
	}

	bodyString := ExecuteRequest(http.NewRequest(http.MethodGet, userUrl, nil))
	expected := "["+
		"{\"name\": \"Timothy\", \"job\": \"Mage\", \"level\": 1}" +
		"{\"name\": \"Shawn\", \"job\": \"Barbarian\", \"level\": 10}" +
		"{\"name\": \"Anthony\", \"job\": \"Hunter\", \"level\": 20}" +
		"]"

	Compare(t, bodyString, expected)
}


func TestCreateInValidUserModels(t *testing.T) {
	ClearAll()
	defer ClearAll()

	cases := []UserModel {
		UserModel{"Anthony", "Programmer", 3},
		UserModel{"Anthony", "Hunter", 0},
		UserModel{"Anthony", "Hunter", 21},
	}

	for _, c := range cases {
		body := CreateUserModelBody(c)
		ExecuteRequest(http.NewRequest(http.MethodPost, userUrl, body))
	}

	bodyString := ExecuteRequest(http.NewRequest(http.MethodGet, userUrl, nil))
	expected := "[]"

	Compare(t, bodyString, expected)
}


// level(>= 1, <= 20), strength(>= 0), dexterity(>= 0), intelligence(>= 0), vitality(>= 0)
// Sum of strength, dexterity, intelligence and vitality values can't exceed level*4
func TestCreateValidItemModels(t *testing.T) {
	ClearAll()
	defer ClearAll()

	cases := []ItemModel {
		ItemModel{"Ring1", 1, 0, 0, 0, 0},
		ItemModel{"Ring2", 10, 2, 2, 2, 2},
		ItemModel{"Ring3", 20, 2, 2, 2, 2},
	}

	for _, c := range cases {
		body := CreateItemModelBody(c)
		ExecuteRequest(http.NewRequest(http.MethodPost, itemUrl, body))
	}

	bodyString := ExecuteRequest(http.NewRequest(http.MethodGet, itemUrl, nil))
	expected := "["+
		"{\"name\": \"Ring1\", \"level\": 1, \"strength\": 1, \"dexterity\": 1, \"intelligence\": 1, \"vitality\": 1}" +
		"{\"name\": \"Ring2\", \"level\": 10, \"strength\": 2, \"dexterity\": 2, \"intelligence\": 2, \"vitality\": 2}" +
		"{\"name\": \"Ring3\", \"level\": 20, \"strength\": 2, \"dexterity\": 2, \"intelligence\": 2, \"vitality\": 2}" +
		"]"

	Compare(t, bodyString, expected)
}


func TestCreateInValidItemModels(t *testing.T) {
	ClearAll()
	defer ClearAll()

	cases := []ItemModel {
		ItemModel{"Ring1", 0, 1, 1, 1, 1},
		ItemModel{"Ring2", 21, 2, 2, 2, 2},
		ItemModel{"Ring3", 20, -1, 2, 2, 2},
		ItemModel{"Ring4", 20, 2, -1, 2, 2},
		ItemModel{"Ring5", 20, 2, 2, -1, 2},
		ItemModel{"Ring6", 20, 2, 2, 2, -1},
		ItemModel{"Ring7", 1, 1, 2, 1, 1},
		ItemModel{"Ring8", 20, 20, 20, 20, 21},
	}

	for _, c := range cases {
		body := CreateItemModelBody(c)
		ExecuteRequest(http.NewRequest(http.MethodPost, itemUrl, body))
	}

	bodyString := ExecuteRequest(http.NewRequest(http.MethodGet, itemUrl, nil))
	expected := "[]"

	Compare(t, bodyString, expected)
}


func ExecuteRequest(req *http.Request, err error) string {
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}


func Compare(t *testing.T, actual string, expected string){
	if actual != expected {
		t.Errorf("response is (%q), instead of %q", actual, expected)
	}
}


func ClearAll(){
	ExecuteRequest(http.NewRequest(http.MethodDelete, deleteUrl, nil))
}


func CreateUserModelBody(user UserModel) io.Reader{

	form := url.Values{}
	form.Add("Name", user.Name)
	form.Add("Job", user.Job)
	form.Add("Level", fmt.Sprintf(":%d", user.Level))

	return strings.NewReader(form.Encode())
}


func CreateItemModelBody(item ItemModel) io.Reader{

	form := url.Values{}
	form.Add("Name", item.Name)
	form.Add("Level", fmt.Sprintf(":%d", item.Level))
	form.Add("Strength", fmt.Sprintf(":%d", item.Strength))
	form.Add("Dexterity", fmt.Sprintf(":%d", item.Dexterity))
	form.Add("Intelligence", fmt.Sprintf(":%d", item.Intelligence))
	form.Add("Vitality", fmt.Sprintf(":%d", item.Vitality))

	return strings.NewReader(form.Encode())
}
