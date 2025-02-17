package semChoose

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

type Assessment struct {
	name string
	path string
}

func SemChoose(url string) {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	if err := spinner.New().Title("Fetching assessments").Action(action).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	params_url := url

	assessments, err := requestClient.SemChooseReq(url)
	if err != nil {
		fmt.Println(helpers.ErrorStyle.Render(fmt.Sprintf("Error: %v\n", err)))
		return
	}

	var selectedOption string
	var assessList []Assessment
	var options []huh.Option[string]

	// Convert assessments to huh options.
	for _, assessment := range assessments {
		assess := Assessment{assessment.Name, assessment.Path}
		assessList = append(assessList, assess)
		options = append(options, huh.NewOption(assess.name, assess.name))
	}
	// Add back option.
	options = append(options, huh.NewOption("Back", "Back"))

	// Create the form.
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Assessments").
				Options(options...).
				Value(&selectedOption),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Handle selection.
	if selectedOption == "Back" {
		inter.UseSemTable(stack.STACK.Pop())
		return
	}

	// Find selected assessment and process it.
	for _, assess := range assessList {
		if assess.name == selectedOption {
			url := configs.BASE_URL + assess.path
			inter.UseYear(url)
			break
		}
	}

	stack.STACK.Push(params_url)
}
