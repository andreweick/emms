// main.go
package main

import (
	"log"
	"os"
	"strings"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

var (
	secret   = os.Getenv("FAUNADB_SECRET")
	endpoint = f.Endpoint("https://db.us.fauna.com")

	client = f.NewFaunaClient(secret, endpoint)

	// dbName = "emms"
)

/*
 * Check for the existencse of the database.  If it exists, return true,
 * otherwise create it.
 */
// func createDatabase() {
// 	res, err := client.Query(
// 		f.If(
// 			f.Exists(f.Database(dbName)),
// 			true,
// 			f.CreateDatabase((f.Obj{"name": dbName}))))

// 	if err != nil {
// 		panic(err)
// 	}

// 	if res != f.BooleanV(true) {
// 		log.Printf("Created Database: %s\n", dbName)
// 	} else {
// 		log.Printf("Database: %s, Already Exists\n", dbName)
// 	}
// }

// func getDbClient() (dbClient *f.FaunaClient) {
// 	var res f.Value
// 	var err error
// 	var dbSecret string

// 	res, err = client.Query(
// 		f.CreateKey(f.Obj{
// 			"database": f.Database(dbName),
// 			"role":     "server",
// 		}))

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = res.At(f.ObjKey("secret")).Get(&dbSecret)

// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Printf("Database: %s, specific key: %s\n", dbName, dbSecret)

// 	dbClient = client.NewSessionClient(dbSecret)

// 	return
// }

func createCollection(collectionName string, dbClient *f.FaunaClient) {

	res, err := dbClient.Query(
		f.If(
			f.Exists(f.Collection(collectionName)),
			true,
			f.CreateCollection(f.Obj{"name": collectionName})))

	if err != nil {
		panic(err)
	}

	if res != f.BooleanV(true) {
		log.Printf("Created Collection: %s\n", collectionName)
	} else {
		log.Printf("Collection: %s, Already Exists\n", collectionName)
	}
}

/*

CreateIndex({
  name: "photographs_search_by_name",
  source: Collection("photographs"),
  terms: [{ field: ["data", "Name"] }]
})

CreateIndex({
  name: "photographs_sort_by_capturedate_asc",
  source: Collection("photographs"),
  values: [
    { field: ["data", "CaptureYearMonthDay"] },
    { field: ["data", "Name"] },
    { field: ["ref"] }
  ]
})

CreateIndex({
  name: "photographs_sort_by_capturedate_desc",
  source: Collection("photographs"),
  values: [
    { field: ["data", "CaptureYearMonthDay"], reverse: true },
    { field: ["data", "Name"] },
    { field: ["ref"] }
  ]
})

## QUERY
Map(
  Paginate(Match(Index("photographs_sort_by_capturedate_asc"))),
  Lambda("pr", Get(Select([2], Var("pr"))))
)
*/

func main() {

	// createDatabase()
	// dbClient := getDbClient()
	createCollection("photographs", client)

	// enumerate directory
	// files, err := os.ReadDir("/Volumes/Mini Pudge/edc/photographs/")
	const sourceDirectory = "./testdata/"
	files, err := os.ReadDir(sourceDirectory)
	if err != nil {
		panic(err)
	}

	for _, fi := range files {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "jpg") {
			pmd := populatePMD(sourceDirectory + fi.Name())

			// parse out the first part of the filename
			parts := strings.Split(fi.Name(), "-")
			pmd.PrefixName = parts[0]

			pmd.Name = fi.Name()
			_, err := client.Query(
				f.Create(
					f.Collection("photographs"),
					f.Obj{"data": pmd},
				),
			)

			if err != nil {
				panic(err)
			}
		}
	}
}
