package main_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	domain "github.com/wheatandcat/PeperomiaBackend/domain"
	"google.golang.org/api/option"
)

var testUID = "5mX0FN6XqVXTrzrlQRMssedkI9R2"

func TestMain(t *testing.T) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	us, _ := getUsers2(ctx, f)
	fmt.Printf("%+v", us)
	cs, _ := getCalendars(ctx, f)
	fmt.Printf("%+v", cs)
	is, _ := getItems(ctx, f)
	fmt.Printf("%+v", is)

}

func getUsers(ctx context.Context, f *firestore.Client) ([]domain.UserRecord, error) {
	var users []domain.UserRecord
	matchItem := f.Collection("version/1/users").Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return users, err
	}

	for _, doc := range docs {
		var user domain.UserRecord
		doc.DataTo(&user)
		users = append(users, user)
	}

	return users, nil
}

func getUsers2(ctx context.Context, f *firestore.Client) ([]domain.UserRecord, error) {
	var users []domain.UserRecord
	matchItem := f.Collection("version/1/users").Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return users, err
	}

	for _, doc := range docs {
		var user domain.UserRecord
		doc.DataTo(&user)
		users = append(users, user)
	}

	return users, nil
}

func getCalendars(ctx context.Context, f *firestore.Client) ([]domain.CalendarRecord, error) {
	var items []domain.CalendarRecord
	matchItem := f.Collection("version/1/users").Doc(testUID).Collection("calendars").Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return items, err
	}

	for _, doc := range docs {
		var item domain.CalendarRecord
		doc.DataTo(&item)
		items = append(items, item)
	}

	return items, nil
}

func getItems(ctx context.Context, f *firestore.Client) ([]domain.ItemRecord, error) {
	var items []domain.ItemRecord
	matchItem := f.CollectionGroup("items").Where("uid", "==", testUID).OrderBy("createdAt", firestore.Desc).Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return items, err
	}

	for _, doc := range docs {
		var item domain.ItemRecord
		doc.DataTo(&item)
		items = append(items, item)
	}

	return items, nil
}
