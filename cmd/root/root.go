package root

import (
	"fmt"
	"os"
	"time"

	"amrita_pyq/cmd/configs"
	"amrita_pyq/cmd/helpers"
	"amrita_pyq/cmd/logo"
	"amrita_pyq/cmd/requestClient"
	"amrita_pyq/cmd/semChoose"
	"amrita_pyq/cmd/semTable"
	"amrita_pyq/cmd/year"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

// Used to implement the interface
type UseInterace struct{}

func (u *UseInterace) UseHuhMenuStart() {
	HuhMenuStart()
}

// Using SemTable from semTable package
func (u *UseInterace) UseSemTable(url string) {
	semTable.SemTable(url)
}

// Using SemChoose from semChoose package
func (u *UseInterace) UseSemChoose(url string) {
	semChoose.SemChoose(url)
}

// Using Year from year package
func (u *UseInterace) UseYear(url string) {
	year.Year(url)
}

type Subject struct {
	name string
	path string
}

var RootCmd = &cobra.Command{
	Use:   "ampyq",
	Short: "Amrita PYQ CLI",
	Long:  `A CLI application to access Amrita Repository for previous year question papers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(helpers.LogoStyle.Render(logo.LOGO_ASCII))
		HuhMenuStart()
	},
}

func HuhMenuStart() {
	action := func() {
		time.Sleep(2 * time.Second)
	}
	if err := spinner.New().Title("Fetching Courses").Action(action).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resources, err := requestClient.GetCoursesReq(configs.COURSE_LIST_URL)
	if err != nil {
		fmt.Println(helpers.ErrorStyle.Render(fmt.Sprintf("Error: %v\n", err)))
		os.Exit(1)
	}

	var selectedOption string
	var subjects []Subject
	var options []huh.Option[string]

	// Convert courses to huh options.
	for _, res := range resources {
		subject := Subject{res.Name, res.Path}
		subjects = append(subjects, subject)
		options = append(options, huh.NewOption(subject.name, subject.name))
	}
	// Add quit option.
	options = append(options, huh.NewOption("Quit", "Quit"))

	// Create the form.
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Available Courses").
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
	if selectedOption == "Quit" {
		fmt.Print(helpers.FetchStatusStyle.Render("Goodbye!\n"))
		os.Exit(0)
	}

	// Find selected subject and process it.
	for _, subject := range subjects {
		if subject.name == selectedOption {
			url := configs.BASE_URL + subject.path
			semTable.SemTable(url)
			break
		}
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
