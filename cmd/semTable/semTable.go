package semTable

import (
	"fmt"
	"os"
	"time"

	"amrita_pyq/cmd/configs"
	"amrita_pyq/cmd/helpers"
	"amrita_pyq/cmd/interfaces"
	"amrita_pyq/cmd/requestClient"
	"amrita_pyq/cmd/stack"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

// Interface to access functions from root package
var inter interfaces.Interface

func Init(n interfaces.Interface) {
	inter = n
}

type Semester struct {
	name string
	path string
}

func SemTable(url string) {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	if err := spinner.New().Title("Fetching Semesters").Action(action).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	semesters, err := requestClient.SemTableReq(url)
	if err != nil {
		fmt.Print(helpers.ErrorStyle.Render(fmt.Sprintf("Error: %v\n", err)))
		return
	}

	var selectedOption string
	var sems []Semester
	var options []huh.Option[string]

	// Convert semesters to huh options.
	for _, sem := range semesters {
		semester := Semester{sem.Name, sem.Path}
		sems = append(sems, semester)
		options = append(options, huh.NewOption(semester.name, semester.name))
	}
	// Add back option.
	options = append(options, huh.NewOption("Back", "Back"))

	// Create the form.
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Semesters").
				Options(options...).
				Value(&selectedOption),
		),
	)

	stack.STACK.Push(url) // Save current URL to stack.

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Handle selection.
	if selectedOption == "Back" {
		inter.UseHuhMenuStart() // Go back to main menu.
		return
	}

	// Find selected semester and process it.
	for _, sem := range sems {
		if sem.name == selectedOption {
			url := configs.BASE_URL + sem.path
			inter.UseSemChoose(url)
			break
		}
	}
}
