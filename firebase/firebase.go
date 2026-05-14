package firebase

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	"github.com/WeatherGod3218/iat-animals-rewrite/logging"
	"github.com/sirupsen/logrus"
)

var database *db.Client
var ctx = context.Background()

func initFirebase() {
	options := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL"),
	}, options)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "initFirebase"}).Fatal("error initializing firebase")
	}

	database, err = app.Database(ctx)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "initFirebase"}).Fatal("error connecting to firebase database")
	}

	logging.Logger.WithFields(logrus.Fields{"module": "firebase", "method": "initFirebase"}).Info("connected to firebase!")
}
