package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	domain "github.com/wheatandcat/PeperomiaBackend/domain"
	"google.golang.org/api/option"
)

// Use a service account

// DataBase is firestore data
type DataBase struct {
	Users       []domain.UserRecord
	Items       []domain.ItemRecord
	ItemDetails []domain.ItemDetailRecord
	Calendars   []domain.CalendarRecord
	PushToken   []domain.PushTokenRecord
}

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	us, _ := getUsers(ctx, client)
	cs, _ := getCalendars(ctx, client)
	is, _ := getItems(ctx, client)
	ids, _ := getItemDetails(ctx, client)
	pts, _ := getPushTokens(ctx, client)

	var db DataBase
	db.Users = us
	db.Items = is
	db.ItemDetails = ids
	db.Calendars = cs
	db.PushToken = pts

	if err := insertUserItem(ctx, client, db); err != nil {
		log.Fatalln(err)
	}
}

func insertItemDetail(ctx context.Context, f *firestore.Client, db DataBase, uid string, calendarID string, itemID string) error {
	var items = f.Collection("version/1/users").Doc(uid).Collection("calendars").Doc(calendarID).Collection("items").Doc(itemID)

	for _, i := range db.ItemDetails {
		if uid == i.UID && itemID == i.ItemID {
			if _, err := items.Collection("itemDetails").Doc(i.ID).Set(ctx, i); err != nil {
				return err
			}
		}
	}

	return nil
}

func insertItem(ctx context.Context, f *firestore.Client, db DataBase, uid string, calendarID string, itemID string) error {
	var calendar = f.Collection("version/1/users").Doc(uid).Collection("calendars").Doc(calendarID)

	for _, i := range db.Items {
		if uid == i.UID && itemID == i.ID {
			if _, err := calendar.Collection("items").Doc(itemID).Set(ctx, i); err != nil {
				return err
			}
			if err := insertItemDetail(ctx, f, db, uid, calendarID, itemID); err != nil {
				return err
			}
		}
	}

	return nil
}

func insertCalendar(ctx context.Context, f *firestore.Client, db DataBase, uid string) error {
	var user = f.Collection("version/1/users").Doc(uid)

	for _, c := range db.Calendars {
		if uid == c.UID {
			if _, err := user.Collection("calendars").Doc(c.ID).Set(ctx, c); err != nil {
				return err
			}
			if err := insertItem(ctx, f, db, uid, c.ID, c.ItemID); err != nil {
				return err
			}
		}
	}

	for _, pt := range db.PushToken {
		if _, err := user.Collection("expoPushTokens").Doc(pt.ID).Set(ctx, pt); err != nil {
			return err
		}
	}

	return nil
}

func insertUserItem(ctx context.Context, f *firestore.Client, db DataBase) error {
	var users = f.Collection("version/1/users")

	for _, u := range db.Users {
		if u.UID != "" {
			if _, err := users.Doc(u.UID).Set(ctx, u); err != nil {
				return err
			}
			if err := insertCalendar(ctx, f, db, u.UID); err != nil {
				return err
			}
		}
	}

	return nil
}

func getUsers(ctx context.Context, f *firestore.Client) ([]domain.UserRecord, error) {
	var users []domain.UserRecord
	matchItem := f.Collection("users").Documents(ctx)
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
	matchItem := f.Collection("calendars").Documents(ctx)
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
	matchItem := f.Collection("items").Documents(ctx)
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

func getItemDetails(ctx context.Context, f *firestore.Client) ([]domain.ItemDetailRecord, error) {
	var items []domain.ItemDetailRecord
	matchItem := f.Collection("itemDetails").Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return items, err
	}

	for _, doc := range docs {
		var item domain.ItemDetailRecord
		doc.DataTo(&item)
		items = append(items, item)
	}

	return items, nil
}

func getPushTokens(ctx context.Context, f *firestore.Client) ([]domain.PushTokenRecord, error) {
	var items []domain.PushTokenRecord
	matchItem := f.Collection("expoPushTokens").Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		return items, err
	}

	for _, doc := range docs {
		var item domain.PushTokenRecord
		doc.DataTo(&item)
		items = append(items, item)
	}

	return items, nil
}
