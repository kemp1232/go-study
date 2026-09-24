package main

import (
	"fmt"
	"strings"
	"time"
)

type Employee struct {
	ID int
	Name string
	Department string
}

type TimeEntry struct {
	employeeID int
	Date time.Time
	Hours float64
}

// use either global var or app struct
// lets try app struct
type App struct {
	employees []Employee
	timeEntries []TimeEntry
}

func main() {
	app := App{
		employees: []Employee{},
		timeEntries: []TimeEntry{},
	}

	/*
		Task: Build a command-line application that keeps employees
		and their time entries in memory. No database and no files yet.

		Functions:
		- Add Employee
		- Add Time Entry
		- List Employees
		- Show Employee Time Entries
		- Calculate Employee Total Hours
		- Calculate Department Total hours
		- Exit menu
	*/

	app.openMenu()
}

func (app *App) openMenu() {
	for {
		fmt.Println()
		fmt.Println("============================================================================")
		fmt.Println()
		fmt.Println("=== Employee Time Tracker ===")
		fmt.Println("1. Add Employee")
		fmt.Println("2. Add Time Entry")
		fmt.Println("3. List Employees")
		fmt.Println("4. Show Employee Time Entries")
		fmt.Println("5. Calculate Employee Total Hours")
		fmt.Println("6. Calculate Department Total Hours")
		fmt.Println("7. Exit")

		fmt.Print("Select an option: ")

		var choice int
		fmt.Scan(&choice)
		
		switch choice {
			case 1:
				app.addEmployee()
			case 2:
				app.addTimeEntry()
			case 3:
				app.listEmployees()
			case 4:
				app.listTimeEntries()
			case 5:
				app.calculateEmployeeTotalHours()
			case 6:
				app.calculateDepartmentTotalHours()
			case 7:
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Option not valid")
		}
	}
}

func (app *App) addEmployee() {
	var employeeName string
	var department string
	fmt.Println("=== Add Employee ===")
	fmt.Print("Employee Name: ")
	fmt.Scan(&employeeName)
	fmt.Println("")
	fmt.Print("Employee Department: ")
	fmt.Scan(&department)

	employeeToBeAdded := Employee{
		ID: len(app.employees) + 1,
		Name: strings.ToTitle(employeeName),
		Department: strings.ToTitle(department),
	}

	app.employees = append(app.employees, employeeToBeAdded)
}

func (app *App) addTimeEntry() {
    var employeeId int
	var hours float64
	fmt.Println("=== Add Time Entry ===")
	fmt.Print("Employee ID: ")
	fmt.Scan(&employeeId)
	fmt.Println("")
	for {
		fmt.Print("hours: ")
		fmt.Scan(&hours)
		if hours <= 0 {
			fmt.Println("Hours must be greater than 0")
			fmt.Print("hours: ")
			fmt.Scan(&hours)
		} else {
			break
		}
	}

	if hours <= 0 {
		fmt.Println("Hours must be greater than 0")
		hours = 0
		fmt.Print("hours: ")
		fmt.Scan(&hours)
	}

	// validate if employee exists
	isEmployeeExists := false
	for _, value := range app.employees {
		if (value.ID == employeeId) {
			isEmployeeExists = true
			break
		}
	}

	if (!isEmployeeExists) {
		fmt.Println("Employee ID does not exists")
		return
	}

	timeEntryToBeAdded := TimeEntry{
		employeeID: employeeId,
		Date: time.Now(),
		Hours: hours,
	}

	app.timeEntries = append(app.timeEntries, timeEntryToBeAdded)
}

func (app *App) listEmployees() {
	fmt.Println("=== List Employee ===")
	for _, value := range app.employees {
		fmt.Printf("%v \t %v \t %v \n", value.ID, value.Name, value.Department)
	}
}

func (app *App) listTimeEntries() {
	fmt.Println("=== List Time Entries ===")
	var employeeId int
	fmt.Print("Employee ID: ")
	fmt.Scan(&employeeId)
	for _, value := range app.timeEntries {
		if (value.employeeID == employeeId) {
			fmt.Printf("%v \t %v \t %v \n", value.employeeID, value.Date, value.Hours)
		}
	}
}

func (app *App) calculateEmployeeTotalHours() {
	var employeeId int
	fmt.Println("=== Calculate Employee Total Hours ===")
	fmt.Print("Employee ID: ")
	fmt.Scan(&employeeId)

	var isEmployeeExists bool = false
	var employee Employee
	for _, value := range app.employees {
		if (value.ID == employeeId) {
			isEmployeeExists = true
			employee = value
			break
		}
	}

	if (!isEmployeeExists) {
		fmt.Println("Employee ID does not exists")
		return
	}

	var totalHours float64 = 0
	for _, value := range app.timeEntries {
		if (value.employeeID == employeeId) {
			totalHours = totalHours + value.Hours
		}
	}

	fmt.Printf("%v's total hours is %v hours", employee.Name, totalHours)
	fmt.Println("")
}

func (app *App) calculateDepartmentTotalHours() {
	var department string
	fmt.Println("=== Calculate Department Total Hours ===")
	fmt.Print("Department: ")
	fmt.Scan(&department)

	formattedDepartment := strings.ToTitle(department)
	var totalHours float64 = 0
	for _, timeEntry := range app.timeEntries {
		for _, employee := range app.employees {
			if (timeEntry.employeeID == employee.ID) {
				if (employee.Department == formattedDepartment) {
					totalHours = totalHours + timeEntry.Hours
				}
			}
		}
	}

	fmt.Printf("%v department total hours is %v hours", department, totalHours)
}