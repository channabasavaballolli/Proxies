package main

// import (
// 	"errors"
// 	"fmt"
// )

// func main() {

// 	result1, err1 := divide(10.0, 20.0)
// 	if err1 != nil {
// 		fmt.Println("Error", err1)
// 	} else {
// 		fmt.Println("Result", result1)
// 	}
// 	result2, err2 := divide(10, 0)
// 	if err2 != nil {
// 		fmt.Println("Error", err2)
// 	} else {
// 		fmt.Println("Result", result2)
// 	}

// }
// func divide(n1, n2 float64) (float64, error) {
// 	if n2 == 0 {
// 		return 0.0, errors.New("Cannot divide by zero")
// 	}
// 	return n1 / n2, nil
// }

// type speaker interface { //speaker is a contract
// 	Speak() string //Speak is method
// }
// type Dog struct{}
// type Cat struct{}

// func (d Dog) Speak() string {
// 	return "AOOOOOOOO"
// }
// func (c Cat) Speak() string {
// 	return "MEOWW"
// }
// func MakeItSpeak(s speaker) {
// 	fmt.Println("Animal says:", s.Speak())
// }
// func main() {
// 	d := Dog{}
// 	c := Cat{}
// 	MakeItSpeak(d)
// 	MakeItSpeak(c)
// }

// import "fmt"

// func main() {
// 	s1 := student{
// 		Name:  "Channu",
// 		ID:    2,
// 		Grade: "A",
// 	}
// 	s2 := student{
// 		Name:  "Nikhil",
// 		ID:    4,
// 		Grade: "C",
// 	}
// 	fmt.Println("Student details before struct modifying")
// 	fmt.Println("Name", s1.Name, "\nID", s1.ID, "\nGrade", s1.Grade)
// 	fmt.Println("Name", s2.Name, "\nID", s2.ID, "\nGrade", s2.Grade)

// 	s1.Updatestudent("B")
// 	s2.Updatestudent("D")
// 	fmt.Println("Student details after struct modifying")
// 	fmt.Println("Name", s1.Name, "\nID", s1.ID, "\nGrade", s1.Grade)
// 	fmt.Println("Name", s2.Name, "\nID", s2.ID, "\nGrade", s2.Grade)
// }

// type student struct {
// 	Name  string
// 	ID    int
// 	Grade string
// }

// func (s *student) Updatestudent(newgrade string) {
// 	s.Grade = newgrade
// }

// func main() {
// 	count := counter{value: 10}
// 	count.IncrementValue()
// 	fmt.Println("Value", count.value)
// 	count.IncrementPointer()
// 	fmt.Println("Pointer", count.value)
// }

// type counter struct {
// 	value int
// }

// func (c counter) IncrementValue() {
// 	c.value++
// 	fmt.Println("value inside the function ", c.value)
// }

// func (c *counter) IncrementPointer() {
// 	c.value++
// }

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	num := 58
// 	fmt.Println("Value 1", num)
// 	p := &num
// 	fmt.Println("Address of num", p)
// 	*p = 25
// 	fmt.Println("Value 2", num)
// 	fmt.Println("Address of num", p)

// }

// package main

// import "fmt"

// func main() {
// 	var db []Student

// 	s1 := Student{
// 		Name:  "Channu",
// 		ID:    1,
// 		Grade: "A",
// 	}

// 	s2 := Student{
// 		Name:  "Nirmala",
// 		ID:    2,
// 		Grade: "B",
// 	}
// 	db = append(db, s1)
// 	db = append(db, s2)
// 	for _, s := range db {
// 		fmt.Println("Students details")
// 		fmt.Println("Name", s.Name)
// 		fmt.Println("ID", s.ID)
// 		fmt.Println("Grade", s.Grade)
// 		fmt.Println("================================")
// 	}
// }

// type Student struct {
// 	Name  string
// 	ID    int
// 	Grade string
// }

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	fruits := []string{
// 		"apple",
// 		"banana",
// 		"chikoo",
// 	}
// 	fruits = append(fruits, "mango")
// 	fmt.Println("fruits are: ", fruits)
// 	fmt.Println("Total Fruits", len(fruits))

// 	myFavs := fruits[1:3]
// 	fmt.Println("My favs", myFavs)

// 	ages := map[string]int{
// 		"nanu": 12,
// 		"ninu": 10,
// 	}
// 	fmt.Println(ages)
// 	ages["rahul"] = 22
// 	fmt.Println(ages)
// 	delete(ages, "rahul")
// 	value, exists := ages["rahul"]
// 	if exists {
// 		fmt.Println("rahul is ", value, "years old")
// 	} else {
// 		fmt.Println("rahul not found")
// 	}
// 	E1 := Employee{
// 		ID:         909,
// 		Name:       "Channu",
// 		Department: "IT",
// 	}
// 	fmt.Println("Employee Details", E1)
// 	E1.Department = "Marketing"
// 	fmt.Println("Employee Details", E1)
// }

// type Employee struct {
// 	ID         int
// 	Name       string
// 	Department string
// }
