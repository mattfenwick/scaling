package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
)

func main() {
	user := "postgres"
	pw := "postgres"
	host := "localhost"
	initDbName := "postgres"
	db, err := webserver.InitializeDB(user, pw, host, initDbName)
	utils.DoOrDie(err)

	insert := `insert into documents (parsed, parse_error) values($1, $2)`
	_, err = db.ExecContext(
		context.TODO(),
		insert,
		json.MustMarshalToString([]any{1, 2, 3, "hi", map[string]string{"qrs": "tuv"}}),
		"")
	utils.DoOrDie(err)

	//dbname := "scaling"
	docs, err := webserver.ReadDocuments(context.TODO(), db)
	utils.DoOrDie(err)
	fmt.Printf("docs: %s\n", json.MustMarshalToString(docs))
	for _, doc := range docs {
		bytes, err := base64.StdEncoding.DecodeString(string(doc.Parsed.([]uint8)))
		fmt.Printf("doc? %T\n", doc.Parsed)
		utils.DoOrDie(err)
		fmt.Printf("???? %s\n", bytes)
	}
}
