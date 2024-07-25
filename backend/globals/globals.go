package globals

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var opt = option.WithCredentialsFile(os.Getenv("DM_API_FIREBASE_CREDENTIALS_FILE"))
var app *firebase.App
var myFirestore *firestore.Client
var myAuth *auth.Client
var globalContext context.Context

func GetFirestore() *firestore.Client {
	return myFirestore
}

func GetAuth() *auth.Client {
	return myAuth
}

func GetGlobalContext() context.Context {
	return globalContext
}

func init(){
	globalContext = context.Background()
	var err error
	app, err = firebase.NewApp(globalContext, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app %v\n", err)
	}
	myFirestore, err = app.Firestore(globalContext)
	if err != nil {
		log.Fatalf("error initializing firestore %v\n", err)
	}
	myAuth, err = app.Auth(globalContext)
	if err != nil {
		log.Fatalf("error initializing auth %v\n", err)
	}
}
