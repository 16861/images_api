package store

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

type Repository struct{}

const SERVER = "mongodb://api_user:api123@ds017185.mlab.com:17185/mdb_test"

const USER_COLLECTION = "users"
const SESSION_COLLECTION = "sessions"

const DBNAME = "mdb_test"

func (r Repository) GetUsers() Users {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("GetUsers: failed to connect to DB, err:", err)
	}

	defer session.Close()

	users := Users{}

	c := session.DB(DBNAME).C(USER_COLLECTION)

	if err := c.Find(nil).All(&users); err != nil {
		fmt.Println("GetUsers: Failed to fetch users:", err)
	}

	return users

}

func (r Repository) AddUser(user User) bool {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("AddUser: failed to connect to DB, err:", err)
	}

	defer session.Close()

	users := r.GetUsers()
	for _, existUser := range users {
		if existUser.Name == user.Name {
			log.Printf("AddUser: user %s exists", user.Name)
			return false
		}
	}
	session.DB(DBNAME).C(USER_COLLECTION).Insert(user)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func (r Repository) getSessions() Sessions {
	session, err := mgo.Dial(SERVER)
	var sessions Sessions

	if err != nil {
		fmt.Println("getSessions: failed to connect to DB, err:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(SESSION_COLLECTION)

	if err := c.Find(nil).All(&sessions); err != nil {
		fmt.Println("getSessions: Failed to fetch sessions:", err)
	}
	return sessions
}

func (r Repository) SetSession(s Session) bool {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("failed to connect to DB, err:", err)
		return false
	}

	defer session.Close()

	ss := r.getSessions()
	for _, es := range ss {
		if es.Name == s.Name {
			log.Printf("SetSession: session for %s is already exists\n", s.Name)
			if es.Token == s.Token {
				return true
			}

			session.DB(DBNAME).C(SESSION_COLLECTION).Update(es, s)
			return true
		}
	}

	session.DB(DBNAME).C(SESSION_COLLECTION).Insert(s)
	return true
}

func (r Repository) CheckIfSessionIsExists(u User) bool {
	ss := r.getSessions()
	for _, es := range ss {
		if es.Name == u.Name {
			return true
		}
	}
	return false
}

func (r Repository) CheckCredential(user User) bool {
	users := r.GetUsers()

	for _, existUser := range users {
		if existUser.Name == user.Name && existUser.Pass == user.Pass {
			return true
		}
	}

	return false
}

func (r Repository) CheckIfTokenIsValid(token string) bool {
	ss := r.getSessions()
	for _, es := range ss {
		if es.Token == token {
			return true
		}
	}
	return false
}
