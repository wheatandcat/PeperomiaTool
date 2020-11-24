package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	domain "github.com/wheatandcat/PeperomiaBackend/domain"
	"google.golang.org/api/option"
)

const location = "Asia/Tokyo"

// Dictionary Dictionaryのタイプ
type Dictionary struct {
	Text    string   `json:"text"`
	Bigrams []string `json:"bigrams"`
}

// Data Dataのタイプ
type Data struct {
	Text []string `json:"text"`
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

	deleteDictionary(ctx, client)

	list, err := getVersion(ctx, client)

	var data Data
	request, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(request, &data)

	for _, text := range data.Text {
		if !contains(list, text) {
			list = append(list, text)
		}
	}

	dictionaryList := []Dictionary{}

	for _, text := range list {
		bigrams, err := ngram(text, 2)
		if err != nil {
			log.Fatalln(err)
		}

		dictionary := Dictionary{
			Text:    text,
			Bigrams: bigrams,
		}
		dictionaryList = append(dictionaryList, dictionary)

	}

	file, _ := json.MarshalIndent(dictionaryList, "", " ")
	_ = ioutil.WriteFile("dictionary.json", file, 0644)

}

func deleteDictionary(ctx context.Context, f *firestore.Client) error {
	batch := f.Batch()

	var dictionary = f.Collection("version/1/dictionary").Documents(ctx)
	docs, err := dictionary.GetAll()
	if err != nil {
		return err
	}

	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	_, err = batch.Commit(ctx)

	return err
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func ngram(targetText string, n int) ([]string, error) {
	sepText := strings.Split(targetText, "")
	var ngrams []string

	if len(sepText) < n {
		r := []string{}
		r = append(r, targetText)
		return r, nil
	}

	for i := 0; i < (len(sepText) - n + 1); i++ {
		ngrams = append(ngrams, strings.Join(sepText[i:i+n], ""))
	}
	return ngrams, nil
}

func getVersion(ctx context.Context, f *firestore.Client) ([]string, error) {
	list := []string{}

	var users = f.Collection("version/1/users").Documents(ctx)
	userDocs, err := users.GetAll()
	if err != nil {
		return list, err
	}

	for _, userDoc := range userDocs {

		var calendars = userDoc.Ref.Collection("calendars").Documents(ctx)
		calendarDocs, err := calendars.GetAll()
		if err != nil {
			return list, err
		}
		for _, calendarDoc := range calendarDocs {

			var items = calendarDoc.Ref.Collection("items").Documents(ctx)
			itemDocs, err := items.GetAll()
			if err != nil {
				return list, err
			}
			for _, itemDoc := range itemDocs {

				var itemDetails = itemDoc.Ref.Collection("itemDetails").Documents(ctx)
				itemDetailDocs, err := itemDetails.GetAll()
				if err != nil {
					return list, err
				}

				for _, itemDetailDoc := range itemDetailDocs {

					var id *domain.ItemDetailRecord
					itemDetailDoc.DataTo(&id)
					title1 := strings.Replace(id.Title, "　", "", -1)
					title := strings.Replace(title1, " ", "", -1)

					if !contains(list, title) {
						list = append(list, title)
					}
				}
			}
		}
	}

	return list, err
}
