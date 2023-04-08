package utils

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginStore struct {
	db map[string]string
}

type TokenSessions struct {
	tokens            map[string]time.Time
	expirationEnabled bool
	expirationTime    time.Duration
	mtx               sync.RWMutex
}

func NewTokenSessions() *TokenSessions {
	ts := new(TokenSessions)
	ts.tokens = make(map[string]time.Time)
	ts.expirationEnabled = true
	ts.expirationTime = time.Duration(60 * time.Minute)
	return ts
}

func (ts *TokenSessions) SetValidity(validity time.Duration) {
  ts.expirationTime = validity
}

func NewLoginStore() *LoginStore {
	ls := new(LoginStore)
	ls.db = make(map[string]string)
	return ls
}

func (loginDB *LoginStore) CheckCredentials(user, password string) bool {
	if loginDB.db == nil {
		return false
	}

	password_sha256, found := loginDB.db[user]
	hasher := sha256.New()
	hasher.Write([]byte(password))
	sha := hex.EncodeToString(hasher.Sum(nil))

	if found {
		if password_sha256 == sha {
			log.Println(user, "successfully logged in")
			return true
		}
	}
	log.Println("login failed for user", user)
	return false

}

func (loginDB *LoginStore) ReadPasswordsFile(passwordFile string) error {
	log.Println("Reading credentials from", passwordFile)
	file, err := os.Open(passwordFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		user_pass := strings.FieldsFunc(line, func(r rune) bool {
			if r == ',' {
				return true
			}
			return false
		})

		log.Println(user_pass)

		if len(user_pass) == 2 {
			if loginDB.db == nil {
				loginDB.db = make(map[string]string)
			}

			loginDB.db[user_pass[0]] = strings.Trim(user_pass[1], "\n")
		} else {
			log.Println("Skipping", user_pass)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func (ts *TokenSessions) StartMaintenance(runInterval time.Duration) {

	if !ts.expirationEnabled {
		return
	}

	ticker := time.NewTicker(runInterval)

	go func() {
		for range ticker.C {

			log.Println("Session keys", len(ts.tokens))
			newTokenMap := make(map[string]time.Time)

			ts.mtx.Lock()
			evictionCounter := 0
			for k, v := range ts.tokens {
				if time.Now().Sub(v) < ts.expirationTime {
					newTokenMap[k] = v
				} else {
					evictionCounter += 1
				}
			}
			ts.mtx.Unlock()
			log.Println("Evicted keys", evictionCounter)
			ts.tokens = newTokenMap
		}
	}()

}

func (ts *TokenSessions) Add() string {
	newToken := GenerateSecureToken(16)
	defer ts.mtx.Unlock()
	ts.mtx.Lock()
	ts.tokens[newToken] = time.Now()
	return newToken
}

func (ts *TokenSessions) Search(token string) (bool, error) {
	defer ts.mtx.RUnlock()
	ts.mtx.RLock()
	tokenCreationTime, found := ts.tokens[token]

	if found {
		if ts.expirationEnabled && (time.Now().Sub(tokenCreationTime) > ts.expirationTime) {
			return false, errors.New("token is expired")
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("token invalid")
	}

}

// TODO func rate limit?

func SimpleAuth(ts *TokenSessions, enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enabled {
			// before request
			tokenFound, err := ts.Search(c.GetHeader("token"))
			if !tokenFound {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}
		}
		c.Next()

	}
}
