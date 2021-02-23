package util

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// firebase接続
func initFiirestore() (*firestore.Client, error) {
	ctx := context.Background()

	root, err := os.Getwd()

	// ディレクトリチェック
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Firestore接続処理
	sa := option.WithCredentialsFile(root + "/data_secret.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

// firestoreにデータを追加
func AddFirestore(doc string, data map[string]interface{}) error {
	ctx := context.Background()

	// 初期設定
	client, errInit := initFiirestore()

	// 初期値の処理にミスがないか
	if errInit != nil {
		return errInit
	}

	// データを追加
	_, _, err := client.Collection("dev").Doc(doc).Collection("data").Add(ctx, data)

	if err != nil {
		return err
	}

	return nil
}

// firestoreのデータを取得
func GetFirestore(doc string) ([]map[string]interface{}, error) {
	ctx := context.Background()

	// 初期設定
	client, errInit := initFiirestore()

	// 初期値の処理にミスがないか
	if errInit != nil {
		return nil, errInit
	}

	// データを追加
	iter := client.Collection("dev").Doc("test").Collection("data").Documents(ctx)

	// docのデータを入れる箱
	var data []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		data = append(data, doc.Data())
	}

	return data, nil
}
