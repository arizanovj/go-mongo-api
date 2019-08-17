package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	ID        primitive.ObjectID
	SecretKey string
	Env       *env.Env
}

func (b *Bearer) tokenValid() bool {
	var bearer Bearer
	collection := b.Env.MDB.Database("logs").Collection("log")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//	value := collection.FindOne(ctx, bson.M{"_id": b.ID})
	err := collection.FindOne(ctx, bson.M{"_id": b.ID}).Decode(&bearer)
	if err == nil {
		return true
	}
	return false
}

//Validate validates the api key against the hash in DB
func (b *Bearer) Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	authToken, _ := authtoken.FromRequest(r)

	s := strings.Split(authToken, "_")
	key, _ := primitive.ObjectIDFromHex(s[0])
	b.ID, b.SecretKey = key, s[1]

	if b.tokenValid() {
		next(w, r)
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "Unauthorized access to this resource")
}

//GenerateHASHFromKey generates hashed version of the api key to be used
func (b *Bearer) GenerateHASHFromKey(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) string {
	h := hmac.New(sha256.New, []byte([]byte(b.SecretKey)))
	return hex.EncodeToString(h.Sum(nil))
}
