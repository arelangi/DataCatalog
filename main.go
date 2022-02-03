package main

func main() {

	a := App{}

	dbName := "datacatalog"
	dbUser := "dcadmin"
	dbPassword := "dczendata"
	a.Initialize(dbName, dbUser, dbPassword)
	a.Run(":3000")

}
