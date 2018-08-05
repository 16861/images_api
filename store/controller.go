package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/dgrijalva/jwt-go"
)

type Controller struct {
	Repository Repository
}

func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Print("GetToken: fail to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !c.Repository.CheckCredential(user) {
		log.Println("GetToken: wrong user credentials for user " + user.Name)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ss := c.Repository.getSessions()
	for _, es := range ss {
		if es.Name == user.Name {
			json.NewEncoder(w).Encode(JwtToken{Token: es.Token, UserName: user.Name})
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Name,
		"password": user.Pass,
	})

	log.Println("Username: " + user.Name)
	log.Println("Password: " + user.Pass)

	tokenString, err := token.SignedString([]byte("secret1"))
	if err != nil {
		fmt.Println("GetToken: error trying to obtain a token")
	}

	c.Repository.SetSession(Session{Name: user.Name, Token: tokenString})
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString, UserName: user.Name})

}

func (c *Controller) ConvertImage(w http.ResponseWriter, r *http.Request) {
	var img ImageRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("Error ConverImage ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(body, &img); err != nil {
		w.WriteHeader(422)
		log.Println(err)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error ConvertImage unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if !c.Repository.CheckIfTokenIsValid(img.Token) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var respPayload ImagesResponse
	for _, imgStr := range img.Images {
		cmd := exec.Command("python3", "main.py")
		ioutil.WriteFile("imgs/"+imgStr.FileName, imgStr.Image, 0644)
		if err := cmd.Run(); err != nil {
			log.Printf("ConvertImage: error while executing python command, err: %s,\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		os.Remove("imgs/" + imgStr.FileName)
		b, _ := ioutil.ReadFile("imgs/" + "COMPRESSED_" + imgStr.FileName)
		imgB64 := base64.StdEncoding.EncodeToString(b)

		respPayload = append(respPayload, ImageStructResponse{
			FileName: "COMPRESSED_" + imgStr.FileName,
			Image:    imgB64,
		})

		os.Remove("imgs/" + "COMPRESSED_" + imgStr.FileName)
		cmd.Process.Kill()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respPayload)
	return
}

func (c *Controller) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error AddUser ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddUser ", err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(422)
		log.Println(err)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error ConvertImage unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if c.Repository.AddUser(user) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
	return
}
