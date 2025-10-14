package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	ID          int
	Fname       string
	Lname       string
	DateOfBirth time.Time
	Email       string
	Address     string
	Gender      string
}

var db *sql.DB

func main() {
	var err error
	db, err = dbSqlConnect()
	if err != nil {
		log.Fatal(err)
	}

	//adding student record to table

	s := Student{
		Fname:       "Guru",
		Lname:       "Raghav",
		DateOfBirth: time.Date(1994, time.August, 14, 23, 51, 42, 0, time.UTC),
		Email:       "guru.raghav@yahoo.com",
		Address:     "Bangalore",
		Gender:      "Male",
	}

	fmt.Println("Adding a new student")
	sID, err := addStudent(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("addSudent id: %v \n", sID)

	/*
		fmt.Scanln()
		fmt.Println("Updating an existing student")
		rowsAffected, err := updateStudent()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("student updated, rows affected : %d\n", rowsAffected)
	*/

	//selecting student by ID
	/* fmt.Println("Retrieving one student")
	fmt.Scanln()

	st, err := studentByID(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("studentByID id: %v \n", st)
	*/
	fmt.Println("Retrieving all the students.. Hit ENTER...")
	fmt.Scanln()
	students, err := fetchStudents()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("fetchStudents count: %v \n", len(students))
	for _, stud := range students {
		fmt.Printf("%+v\n", stud)
	}
}

func dbSqlConnect() (*sql.DB, error) {
	// Opening a database connection.
	db, err := sql.Open("mysql", "root:rootuser@tcp(localhost:3306)/go_db_demo?parseTime=true")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!")
	return db, nil
}

func addStudent(s Student) (int64, error) {
	query := "insert into students (fname, lname, date_of_birth, email, gender, address) values (?, ?, ?, ?, ?, ?);"
	result, err := db.Exec(query, s.Fname, s.Lname, s.DateOfBirth, s.Email, s.Gender, s.Address)
	if err != nil {
		return 0, fmt.Errorf("addStudent Error: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addSudent Error: %v", err)
	}

	return id, nil
}

func updateStudent() (int64, error) {
	query := "update students set fname=?, lname=? where ID=?;"
	result, err := db.Exec(query, "Magesh", "Kuppan", 1)
	if err != nil {
		return 0, fmt.Errorf("updateStudent Error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateStudent Error: %v", err)
	}

	return rowsAffected, nil
}

func fetchStudents() ([]Student, error) {
	// A slice of Students to hold data from returned rows.
	var students []Student

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {

		// create an empty student object
		var s Student

		// populate the data from the cursor to the go object
		if err := rows.Scan(&s.ID, &s.Fname, &s.Lname, &s.DateOfBirth, &s.Email, &s.Address, &s.Gender); err != nil {
			return nil, fmt.Errorf("fetchStudents %v", err)
		}

		// append the populated object to the slice
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}

	return students, nil
}

func studentByID(id int64) (Student, error) {
	var st Student

	row := db.QueryRow("SELECT * FROM students WHERE id = ?", id)
	if err := row.Scan(&st.ID, &st.Fname, &st.Lname, &st.DateOfBirth, &st.Email, &st.Address, &st.Gender); err != nil {
		if err == sql.ErrNoRows {
			return st, fmt.Errorf("studentById %d: no such student", id)
		}
		return st, fmt.Errorf("studentById %d: %v", id, err)
	}
	return st, nil
}
