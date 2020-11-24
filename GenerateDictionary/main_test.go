package main_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Item struct {
	Text    string   `json:"text" firestore:"text"`
	Bigrams []string `json:"bigrams" firestore:"bigrams"`
}

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

	searchWords := "公園"

	bigrams, err := ngram(searchWords, 2)
	if err != nil {
		log.Fatalln(err)
	}

	var query = f.Collection("version/1/dictionary").Limit(1000)

	for _, bigram := range bigrams {
		key := "bigrams." + bigram
		query = query.Where(key, "==", true)
	}

	matchItem := query.Documents(ctx)
	docs, err := matchItem.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	var items []Item

	for _, doc := range docs {
		var item Item
		doc.DataTo(&item)
		items = append(items, item)
	}

	fmt.Printf("%+v", items)

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
