package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arizanovj/go-mongo-api/env"
	"github.com/harlow/authtoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Bearer represents the authentication logic for api key authentication through bearer token
type Bearer struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SecretKey string             `json:"api_key,omitempty" bson:"api_key,omitempty"`
	Env       *env.Env           `json:"-"  bson:"-"`
}

var collection = "user"

//validqate the token provided in the bearer header
func (b *Bearer) tokenValid() error {
	var bearer Bearer
	collection := b.Env.MDB.Database(b.Env.DBName).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"_id": b.ID}).Decode(&bearer)
	if err != nil {
		return errors.New("no such API user")
	}
	if b.hashFromString(b.SecretKey) == bearer.SecretKey {
		return nil
	}
	return errors.New("key is not authenticated")
}

//create sha256 hash from string
func (b *Bearer) hashFromString(stringToHash string) string {
	h := hmac.New(sha256.New, []byte([]byte(stringToHash)))
	return hex.EncodeToString(h.Sum(nil))
}

//Validate validates the api key against the hash in DB
func (b *Bearer) Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authToken, err := authtoken.FromRequest(r)
	if err != nil {
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	s := strings.Split(authToken, "_")
	if len(s) != 2 {
		fmt.Fprint(w, "token should be identity key and secret key separated by '_'")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	key, err := primitive.ObjectIDFromHex(s[0])
	if err != nil {
		fmt.Fprint(w, "invalid identity key")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b.ID, b.SecretKey = key, s[1]

	err = b.tokenValid()

	if err == nil {
		next(w, r)
		return
	}
	fmt.Fprint(w, err)
	w.WriteHeader(http.StatusUnauthorized)

}
