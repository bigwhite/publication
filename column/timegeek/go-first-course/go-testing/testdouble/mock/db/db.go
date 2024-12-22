package main

// Define the `Database` interface
type Database interface {
	Save(data string) error
	Get(id int) (string, error)
}

// Example functions that use the `Database` interface
func saveData(db Database, data string) error {
	return db.Save(data)
}

func getData(db Database, id int) (string, error) {
	return db.Get(id)
}
