// main.go
package main

import (
	"log"
	"os"
	"strings"
	"time"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

var (
	secret   = os.Getenv("FAUNADB_SECRET")
	endpoint = f.Endpoint("https://db.us.fauna.com")

	client = f.NewFaunaClient(secret, endpoint)

	dbName = "emms"
)

/*
 * Check for the existencse of the database.  If it exists, return true,
 * otherwise create it.
 */
func createDatabase() {
	res, err := client.Query(
		f.If(
			f.Exists(f.Database(dbName)),
			true,
			f.CreateDatabase((f.Obj{"name": dbName}))))

	if err != nil {
		panic(err)
	}

	if res != f.BooleanV(true) {
		log.Printf("Created Database: %s\n", dbName)
	} else {
		log.Printf("Database: %s, Already Exists\n", dbName)
	}
}

func getDbClient() (dbClient *f.FaunaClient) {
	var res f.Value
	var err error
	var dbSecret string

	res, err = client.Query(
		f.CreateKey(f.Obj{
			"database": f.Database(dbName),
			"role":     "server",
		}))

	if err != nil {
		panic(err)
	}

	err = res.At(f.ObjKey("secret")).Get(&dbSecret)

	if err != nil {
		panic(err)
	}

	log.Printf("Database: %s, specific key: %s\n", dbName, dbSecret)

	dbClient = client.NewSessionClient(dbSecret)

	return
}

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

func main() {

	createDatabase()
	dbClient := getDbClient()
	createCollection("photographs", dbClient)

	pmd := PhotoMetaData{
		Name:                "test.jpg",
		ParsedName:          "testParsedName",
		Artist:              "testArtist",
		CaptureTime:         time.Now(),
		CaptureYear:         "2020",
		CaptureYearMonth:    "2020-01",
		CaptureYearMonthDay: "2020-01-01",
		Description:         "testDescription",
		Caption:             "testCaption",
		ID:                  1,
		Height:              100,
		Width:               100,
	}

	_, err := dbClient.Query(
		f.Create(
			f.Collection("photographs"),
			f.Obj{"data": pmd},
		),
	)

	if err != nil {
		panic(err)
	}

	// enumerate directory
	files, err := os.ReadDir("/Volumes/Mini Pudge/edc/photographs/")
	if err != nil {
		panic(err)
	}

	for _, fi := range files {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "jpg") {
			pmd := populatePMD("/Volumes/Mini Pudge/edc/photographs/" + fi.Name())

			pmd.Name = fi.Name()
			_, err := dbClient.Query(
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

	// res, err := client.Query(f.Get(f.Ref(f.Collection("products"), "202")))
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(res)
}
