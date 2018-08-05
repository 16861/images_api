package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"rest-and-go/store"
)

const HOST = "http://localhost"
const PORT = "3001"
const GET_TOKEN_ROUTING = "/get_token"
const ADD_USER_ROUTING = "/adduser"
const CONVERT_IMAGE_ROUTING = "/convertimage"

func getResponse(route string, data *bytes.Buffer) *http.Response {
	res, _ := http.Post(
		HOST+":"+PORT+route,
		"application/json; charset=utf-8",
		data,
	)
	return res
}

func TestAddUser(t *testing.T) {
	user := store.User{
		Name: "igor",
		Pass: "123",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(user)
	res := getResponse(ADD_USER_ROUTING, b)

	if res == nil {
		t.Errorf("TestAddUser: can't revieve response from server, api is shut down")
		return
	}
	if res.StatusCode == http.StatusOK {
		return
	} else if res.StatusCode == http.StatusCreated {
		return
	}

	t.Errorf("TestAddUser: wrong status code, actual: %d\n", res.StatusCode)
}

func TestGetToken(t *testing.T) {

	user := store.User{
		Name: "igor",
		Pass: "123",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(user)
	res := getResponse(GET_TOKEN_ROUTING, b)

	if res == nil {
		t.Errorf("TestGetToken: can't revieve response from server, api is shut down")
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("TestGetToken: wrong status code, actual: %d\n", res.StatusCode)
		return
	}
	var token store.JwtToken
	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(body, &token)
	if err != nil {
		t.Errorf("TestGetToken: can't umnarshal response from server. err: %s", err)
		return
	}

	if token.UserName != user.Name {
		t.Errorf("TestGetToken: name mismatch. in response: %s", token.UserName)
		return
	}
	if token.Token == "" {
		t.Errorf("TestGetToken: token is empty for %s", token.UserName)
		return
	}
}

func TestGetTokenBadRequest(t *testing.T) {
	b := new(bytes.Buffer)

	res := getResponse(GET_TOKEN_ROUTING, b)

	if res == nil {
		t.Errorf("TestGetTokenBadRequest: can't revieve response from server, api is shut down")
		return
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("TestGetTokenBadRequest: wrong status code, actual: %d\n", res.StatusCode)
		return
	}
	body, _ := ioutil.ReadAll(res.Body)
	if len(body) > 0 {
		t.Errorf("TestGetTokenBadRequest: response body is not empty")
		return
	}
}

func TestGetTokenFordidden(t *testing.T) {
	user := store.User{
		Name: "igor",
		Pass: "123j",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(user)

	res := getResponse(GET_TOKEN_ROUTING, b)

	if res == nil {
		t.Errorf("TestGetTokenFordidden: can't revieve response from server, api is shut down")
		return
	}
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("TestGetTokenFordidden: wrong status code, actual: %d\n", res.StatusCode)
		return
	}
	body, _ := ioutil.ReadAll(res.Body)
	if len(body) > 0 {
		t.Errorf("TestGetTokenFordidden: response body is not empty")
		return
	}
}

func TestConvertImage(t *testing.T) {
	user := store.User{
		Name: "igor",
		Pass: "123",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(user)
	res := getResponse(GET_TOKEN_ROUTING, b)

	var token store.JwtToken
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &token); err != nil {
		t.Errorf("TestConvertImage: failed to unmarshal response fromn server")
		return
	}

	res, err := http.Get("https://golang.org/doc/gopher/frontpage.png")
	if err != nil {
		t.Errorf("TestConvertImage: failed to fetch image")
		return
	}

	defer res.Body.Close()
	imgReq := store.ImageRequest{
		Token: token.Token,
		Images: []store.ImageStruct{
			store.ImageStruct{
				FileName: "gopher.png",
			},
			store.ImageStruct{
				FileName: "clouds-cold.jpg",
			},
		},
	}

	res2, err := http.Get("https://images.pexels.com/photos/772803/pexels-photo-772803.jpeg?cs=srgb&dl=altitude-clouds-cold-772803.jpg&fm=jpg")
	if err != nil {
		t.Errorf("TestConvertImage: failed to fetch image")
		return
	}

	defer res2.Body.Close()
	imgReq.Images[0].Image, _ = ioutil.ReadAll(res.Body)
	imgReq.Images[1].Image, _ = ioutil.ReadAll(res2.Body)

	b = new(bytes.Buffer)
	json.NewEncoder(b).Encode(imgReq)
	res = getResponse(CONVERT_IMAGE_ROUTING, b)

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("TestConvertImage: wrong status code, actual: %d\n", res.StatusCode)
		return
	}

	/* 	respBody, _ := ioutil.ReadAll(res.Body)
	   	ioutil.WriteFile("body", respBody, 0644) */

}
