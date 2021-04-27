package mongodb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func setupDB(t *testing.T) *MongoDB {
	db := MongoDB{
		"mongodb://localhost", // place your mLabs URL here for testing
		"currency",
		"testStudents",
	}

	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()

	if err != nil {
		t.Error(err)
	}
	return &db
}

func tearDownDB(t *testing.T, db *MongoDB) {
	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	if err != nil {
		t.Error(err)
	}

	err = session.DB(db.DatabaseName).DropDatabase()
	if err != nil {
		t.Error(err)
	}
}

// Testing the Upsert statement
func TestMongo_Upsert(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()

	student := Student{Name: "Tom", Age: 21, StudentID: "id1"}

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var resultStudent Student
	upsertInfo, err := session.DB(db.DatabaseName).C(db.CollectionName).Upsert(student, student)

	if err != nil {
		fmt.Printf("error in FindId(): %v", err.Error())
		return
	}

	id := upsertInfo.UpsertedId
	err = session.DB(db.DatabaseName).C(db.CollectionName).FindId(id).One(&resultStudent)

	if err != nil {
		fmt.Printf("error in FindId(): %v", err.Error())
		return
	}

	if db.Count() != 1 {
		t.Error("adding new student failed.")
	}
}

// Testing the Insert and FindId statements
func TestMongoDB_Insert(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized. student count() should be 0.")
		return
	}

	student := Student{Name: "Tom", Age: 21, StudentID: "id1", Id: bson.NewObjectId()}
	db.Add(student)

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var resultStudent Student
	err = session.DB(db.DatabaseName).C(db.CollectionName).FindId(student.Id).One(&resultStudent)

	if err != nil {
		fmt.Printf("error in FindId(): %v", err.Error())
		return
	}

	if db.Count() != 1 {
		t.Error("adding new student failed.")
	}
}

func TestStudentsMongoDB_Get(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized. student count() should be 0.")
		return
	}

	student := Student{Name: "Tom", Age: 21, StudentID: "id1"}
	db.Add(student)

	if db.Count() != 1 {
		t.Error("adding new student failed.")
	}

	newStudent, ok := db.Get(student.StudentID)
	if !ok {
		t.Error("couldn't find Tom")
	}

	if newStudent.Name != student.Name ||
		newStudent.Age != student.Age ||
		newStudent.StudentID != student.StudentID {
		t.Error("students do not match")
	}

	all := db.GetAll()
	if len(all) != 1 || all[0].StudentID != student.StudentID {
		t.Error("GetAll() doesn't return proper slice of all the items")
	}
}
