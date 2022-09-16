package main

import (
	"encoding/json"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

//
// Function:          RenderDialogContents
//
// Description:       This function is used to process and render the contents of a dialog.
//
// Inputs:
//                   template      The template to use
//                   data          The data to use to render the template
//
func RenderDialogContents(template string, data map[string]string) string {
	//
	// Render the current for the first pass.
	//
	page, err := raymond.Render(template, data)
	if err != nil {
		log.Fatal(err)
	}

	//
	// Return the results.
	//
	return page
}

//
// Function:     putRequest
//
// Description:  This method will issue a put request with the data sent
//               as json in the body.
//
// Inputs:
//               url        The url to send the request
//               data       An io.Reader pointing to a json string
//
func putRequest(url string, data io.Reader) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err2 := client.Do(req)
	if err2 != nil {
		// handle error
		log.Fatal(err2)
	}
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		log.Fatal(err3)
	}
	resp.Body.Close()
	return string(body)
}

//
// Function:     fileExists
//
// Description:  This function checks if a file exists and is not a directory before we
//               try using it to prevent further errors.
//
// Inputs:       filename       A string representing the file to check.
//
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//
// Function:     FilenameWithoutExtension
//
// Description:  This function trims the extension off of a file name.
//
// Inputs:
//               fn      File name to remove the extension.
//
func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

//
// Function:     Map
//
// Description:  A utility function to return an array of strings
//               that was processed by a given function.
//
// Inputs:
//               list      Array of strings
//               f         Function to execute on each string
//
func Map(list []string, f func(string) string) []string {
	result := make([]string, len(list))
	for i, item := range list {
		result[i] = f(item)
	}
	return result
}

//
// NOTE: This section is for the build cli using Bubbletea framework.
//
//
// Struct:		model
//
// description: The structure for the bubbletea interface for building a dialog.
//

//
// This is the information that we will be filling up to make the
// dialogs.
//
type DialogItem struct {
	ModelType string `json:"modaltype" bindin:"required"`
	Name      string `json:"name" bindin:"required"`
	Id        string `json:"id" bindin:"required"`
	Value     string `json:"value" bindin:"required"`
	For       string `json:"for" bindin:"required"`
}

type DialogButton struct {
	Name   string `json:"name" bindin:"required"`
	Id     string `json:"id" bindin:"required"`
	Action string `json:"action" bindin:"required"`
}

type ModalDialog struct {
	Items   []DialogItem   `json:"items" bindin:"required"`
	Buttons []DialogButton `json:"buttons" bindin:"required"`
}

var buildDialog ModalDialog // The dialog structure we need to build

type model struct {
	savefile    string            // The file to save the structure
	orgItems    []string          // beginning list of choices
	diagItems   []string          // These are the choices for adding to a dialog
	choices     []string          // These are the current items being shown.
	cursor      int               // which to-do list item our cursor is pointing at
	selected    int               // which to-do items are selected
	state       int               // What state the system is in
	inputs []textinput.Model      // This contains the input fields for the labels
	focused     int               // This is the currently focused input
  currentQueue []int            // The queue of inputs to use
  labelqueue   []int            // The queue of inputs for a label
  inputqueue   []int            // The queue of inputs for a label
  buttonqueue  []int            // The queue of inputs for a button
	err         error             // this will contain any errors from the validators
}

type tickMsg struct{}
type errMsg error

const (
	name = iota
	id
	value
	forid
)

const (
	purple  = lipgloss.Color("#9580FF")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(purple)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

// Validator functions to ensure valid input
func nameValidator(s string) error {
	return nil
}

func stringValidator(s string) error {
	return nil
}

func initialModel(savefile string) model {
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[name] = textinput.New()
	inputs[name].Placeholder = ""
	inputs[name].Focus()
	inputs[name].CharLimit = 50
	inputs[name].Width = 52
	inputs[name].Prompt = ""
	inputs[name].Validate = nameValidator

	inputs[id] = textinput.New()
	inputs[id].Placeholder = ""
	inputs[id].CharLimit = 50
	inputs[id].Width = 52
	inputs[id].Prompt = ""
	inputs[id].Validate = nameValidator

	inputs[value] = textinput.New()
	inputs[value].Placeholder = ""
	inputs[value].CharLimit = 0
	inputs[value].Prompt = ""
	inputs[value].Validate = stringValidator

	inputs[forid] = textinput.New()
	inputs[forid].Placeholder = ""
	inputs[forid].CharLimit = 50
	inputs[forid].Width = 52
	inputs[forid].Prompt = ""
	inputs[forid].Validate = nameValidator

	return model{
		// Our list of acctions
		savefile:    savefile,
		orgItems:    []string{"Add Item", "Add Button", "Save"},
		diagItems:   []string{"Add label", "Add Input", "Save"},
		choices:     []string{"Add Item", "Add Button", "Save"},
		cursor:      0,
		state:       0,
		inputs: inputs,
    currentQueue: []int{name, id, value, forid},
    labelqueue:   []int{name, id, value, forid},
    inputqueue:   []int{name, id, value},
    buttonqueue:  []int{name, id, value, forid},
		focused:     0,
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (m *model) nextInput() {
  //
  // Increment the focused item and wrap around if
  // too large.
  //
	m.focused = (m.focused + 1) % len(m.currentQueue)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
  //
  // Decrement the focused item.
  //
	m.focused--

	//
  // If less than zero, wrap around to the highest number.
  //
	if m.focused < 0 {
    m.focused = len(m.currentQueue) - 1
	}
}

type makeItemFinishedMsg struct{ m model }

func (m model) MakeItem() tea.Msg {
	return makeItemFinishedMsg{m}
}

type makeLabelFinishedMsg struct{ m model }

func (m model) MakeLabel() tea.Msg {
	return makeLabelFinishedMsg{m}
}

type makeButtonFinishedMsg struct{ m model }

func (m model) MakeButton() tea.Msg {
	return makeButtonFinishedMsg{m}
}

type makeInputFinishedMsg struct{ m model }

func (m model) MakeInput() tea.Msg {
	return makeInputFinishedMsg{m}
}

type labelInputFinishedMsg struct{ m model }

func (m model) SaveInput() tea.Msg {
	//
	// Save the input data to the build structure
	//
	switch m.state {
	case 2:
		//
		// This is the label case.
		//
		var di DialogItem
		di.ModelType = "label"
		di.Name = m.inputs[name].Value()
		di.Id = m.inputs[id].Value()
		di.Value = m.inputs[value].Value()
		di.For = m.inputs[forid].Value()
    m.inputs[name].Reset()
    m.inputs[name].SetValue("")
	  m.inputs[name].Focus()
    m.inputs[id].Reset()
    m.inputs[id].SetValue("")
    m.inputs[id].Blur()
    m.inputs[forid].Reset()
    m.inputs[forid].SetValue("")
    m.inputs[forid].Blur()
    m.inputs[value].Reset()
    m.inputs[value].SetValue("")
    m.inputs[value].Blur()
		buildDialog.Items = append(buildDialog.Items, di)
		break

	case 4:
		//
		// Creating a Input
		//
		var di DialogItem
		di.ModelType = "input"
		di.Name = m.inputs[name].Value()
		di.Id = m.inputs[id].Value()
		di.Value = m.inputs[value].Value()
		di.For = ""
    m.inputs[name].Reset()
    m.inputs[name].SetValue("")
	  m.inputs[name].Focus()
    m.inputs[id].Reset()
    m.inputs[id].SetValue("")
    m.inputs[id].Blur()
    m.inputs[forid].Reset()
    m.inputs[forid].SetValue("")
    m.inputs[forid].Blur()
    m.inputs[value].Reset()
    m.inputs[value].SetValue("")
    m.inputs[value].Blur()
		buildDialog.Items = append(buildDialog.Items, di)
		break

	case 6:
		//
		// Creating a button
		//
    var db DialogButton
    db.Name = m.inputs[name].Value()
    db.Id = m.inputs[id].Value()
    db.Action = m.inputs[value].Value()
    m.inputs[name].Reset()
    m.inputs[name].SetValue("")
	  m.inputs[name].Focus()
    m.inputs[id].Reset()
    m.inputs[id].SetValue("")
    m.inputs[id].Blur()
    m.inputs[forid].Reset()
    m.inputs[forid].SetValue("")
    m.inputs[forid].Blur()
    m.inputs[value].Reset()
    m.inputs[value].SetValue("")
    m.inputs[value].Blur()
    buildDialog.Buttons = append(buildDialog.Buttons, db)
    break

	default:
		break
	}

	//
	// Go back to the first state.
	//
	return labelInputFinishedMsg{m}
}

type saveSturctureFinishedMsg struct{ m model }

func (m model) SaveStructure() tea.Msg {
	//
	// Save the structure to a file.
	//
	file, _ := json.MarshalIndent(buildDialog, "", " ")
	header := "# This a dialog created by the builder.\n"
	_ = ioutil.WriteFile(m.savefile, []byte(header+string(file)), 0644)
	return saveSturctureFinishedMsg{m}
}

func switchInQueryMode(m model, msg string) (tea.Model, tea.Cmd) {
	// Cool, what was the actual key pressed?
	switch msg {

	// These keys should exit the program.
	case "ctrl+c", "q":
		return m, tea.Quit

	// The "up" and "k" keys move the cursor up
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	// The "down" and "j" keys move the cursor down
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	// The "enter" key will select the action to perform.
	case "enter":
		switch m.state {
		case 0:
			if m.cursor == 0 {
				return m, m.MakeItem
			} else if m.cursor == 1 {
				return m, m.MakeButton
			} else {
				// this would save.
				return m, m.SaveStructure
			}

		case 1:
			if m.cursor == 0 {
				return m, m.MakeLabel
			} else if m.cursor == 1 {
				return m, m.MakeInput
			} else if m.cursor == 2 {
				// This would save.
				return m, m.SaveStructure
			}

		case 5:
			return m, m.MakeButton
		}
	}
	return m, nil
}

func switchInLabelMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
      max := len(m.inputs) - 1
      if m.state == 4 || m.state == 6 {
        max = max - 1
      }
			if m.focused == max {
				//
				// This is the last input, save the inputs
				//
				return m, m.SaveInput
			} else {
				m.nextInput()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.currentQueue {
			m.inputs[m.currentQueue[i]].Blur()
		}

    //
    // Make sure the focused item is in the current queue.
    //
    if !contains(m.currentQueue, m.focused) {
      //
      // Not there, reset it.
      //
      m.focused = m.currentQueue[0]
    }

    //
    // Focus the current input.
    //
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.currentQueue {
		m.inputs[m.currentQueue[i]], cmds[m.currentQueue[i]] = m.inputs[m.currentQueue[i]].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

//
// Checking for item inside of an array.
//
func contains(a []int, item int) bool {
	for _, v := range a {
		if v == item {
			return true
		}
	}
	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg2 := msg.(type) {

	case makeItemFinishedMsg:
		m.choices = m.diagItems
    m.cursor = 0
		m.state = 1
		return m, nil

	case makeLabelFinishedMsg:
		m.choices = m.orgItems
    m.cursor = 0
		m.state = 2
    m.currentQueue = m.labelqueue
    m.focused = name
		return m, nil

	case labelInputFinishedMsg:
		m.choices = m.orgItems
    m.cursor = 0
		m.state = 0
		return m, nil

	case makeInputFinishedMsg:
		m.choices = m.orgItems
    m.cursor = 0
		m.state = 4
    m.currentQueue = m.inputqueue
    m.focused = name
		return m, nil

	case makeButtonFinishedMsg:
		m.choices = m.orgItems
    m.cursor = 0
		m.state = 6
    m.currentQueue = m.buttonqueue
    m.focused = name
		return m, nil

	case saveSturctureFinishedMsg:
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		switch m.state {
		case 0, 1, 3, 5:
			return switchInQueryMode(m, msg2.String())
		case 2, 4, 6:
			return switchInLabelMode(m, msg)
    }
	}
	return m, nil
}

func viewChoices(m model) string {
	// The header
	s := "\n\n\nWhat do you want to do?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress j to move down. Press k to move up. Press enter to select. Press q to quit.\n\n\n\n"

	// Send the UI for rendering
	return s
}

func viewLabelInputs(m model) string {
	return fmt.Sprintf(
		` Fields for the Label

 %s
 %s
 %s
 %s
 %s  
 %s
 %s  
 %s
 %s
`,
		inputStyle.Width(10).Render("Label Name"),
		m.inputs[name].View(),
		inputStyle.Width(2).Render("ID"),
		m.inputs[id].View(),
		inputStyle.Width(5).Render("Value"),
		m.inputs[value].View(),
		inputStyle.Width(6).Render("For ID"),
		m.inputs[forid].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func viewInputInputs(m model) string {
	return fmt.Sprintf(
		` Fields for the Input

 %s
 %s
 %s
 %s
 %s  
 %s
 %s  
`,
		inputStyle.Width(10).Render("Input Name"),
		m.inputs[name].View(),
		inputStyle.Width(2).Render("ID"),
		m.inputs[id].View(),
		inputStyle.Width(13).Render("Default Value"),
		m.inputs[value].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func viewButtonInputs(m model) string {
	return fmt.Sprintf(
		` Fields for a Button

 %s
 %s
 %s
 %s
 %s  
 %s
 %s  
`,
		inputStyle.Width(11).Render("Button Name"),
		m.inputs[name].View(),
		inputStyle.Width(2).Render("ID"),
		m.inputs[id].View(),
		inputStyle.Width(13).Render("Action"),
		m.inputs[value].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}


//
// Function:    View
//
// Description: The view on a model controls how it is displayed. It returns strings
//              for displaying to the user. We select the right view based on the
//              state of the statemachine.
//
func (m model) View() string {
	switch m.state {
	case 0, 1, 3, 5:
		return viewChoices(m)
	case 2:
		return viewLabelInputs(m)
  case 4:
	  return viewInputInputs(m)
  case 6:
    return viewButtonInputs(m)
	}
  return viewChoices(m)
}

//
// Function:     main
//
// Description:  This is the main entry point for the program.
//
// Inputs:
//               The inputs are assigned to os.Argv. It should be a dialog
//               name and the data to use to expand it. Currently made for the
//				 macOS.
//
// #TODO: make to work for other Oses.
//
func main() {
	dialog := ""
	data := make(map[string]string, 0)
	if len(os.Args) > 1 {
		//
		// Get the command or dialog to process.
		//
		dialog = os.Args[1]

		//
		// Get the two template locations.
		//
		home := os.Getenv("HOME")
		progHome, _ := os.Executable()
		progHome = filepath.Dir(progHome)
		templates1 := filepath.Join(progHome, "../Resources/dialogs")
		templates2 := filepath.Join(home, ".config/scriptserver/dialogs")

		//
		// Process the command or the template.
		//
		switch dialog {
		case "list":
			{
				//
				// Give the user a json list of dialogs in the program
				// area and in the user directory.
				//
				var nlist []string
				file, _ := os.Open(templates1)
				dlist, _ := file.Readdirnames(0) // 0 to read all files and folders
				for _, name := range dlist {
					nlist = append(nlist, name)
				}
				file.Close()
				file, _ = os.Open(templates2)
				dlist, _ = file.Readdirnames(0)
				for _, name := range dlist {
					nlist = append(nlist, name)
				}
				file.Close()
				nlist = Map(nlist, FilenameWithoutExtension)
				pjson, err := json.Marshal(nlist)
				if err != nil {
					log.Fatal("Cannot encode to JSON ", err)
				}
				fmt.Printf("{ \"dialogs\": %s}\n", pjson)
			}
		case "build":
			{
				//
				// We are going to build a dialog.
				//
				if len(os.Args) < 3 {
					fmt.Print("\nNot enough arguments. Please give the the dialog a name!\n")
				} else {
					//
					// Initialize the buildDialog  structure. I don't have the buttons done yet, but to test
					// what I do have has to have this structure. But, every dialog needs a cancel button.
					//
					buildDialog.Buttons = make([]DialogButton, 1)
					buildDialog.Buttons[0].Name = "Cancel"
					buildDialog.Buttons[0].Id = "cancel"
					buildDialog.Buttons[0].Action = "cancel"

					//
					// create the Bubbletea interface for building the new dialog
					//
					p := tea.NewProgram(initialModel(filepath.Join(templates2, fmt.Sprintf("%s%s", os.Args[2], ".json"))))
					if err := p.Start(); err != nil {
						fmt.Printf("Alas, there's been an error: %v", err)
						os.Exit(1)
					}
				}
			}
		default:
			{
				//
				// Create the rest of the command line into the data needed for the dialog template.
				//
				for i := 2; i < len(os.Args); i++ {
					data[fmt.Sprintf("data%d", i-1)] = os.Args[i]
				}

				//
				// Create an error dialog if the dialog can't be found.
				//
				var jsonStr string = "{ \"html\": \"<h1>Dialog not found.<h1>\", \"width\": 100, \"height\": 200, \"x\": 200, \"y\": 200}"

				//
				// Create the two possible file locations.
				//
				templatefile1 := filepath.Join(templates1, fmt.Sprintf("%s.json", dialog))
				templatefile2 := filepath.Join(templates2, fmt.Sprintf("%s.json", dialog))
				if fileExists(templatefile1) {
					//
					// The dialog is in the Resources directory of the application bundle
					//
					Str, _ := ioutil.ReadFile(templatefile1)
					jsonStr = string(Str)
				} else if fileExists(templatefile2) {
					//
					// The dialog is in the user's home directory area.
					//
					Str, _ := ioutil.ReadFile(templatefile2)
					jsonStr = string(Str)
				}
				if jsonStr[0] == '#' {
					//
					// This is a dialog build using a json structure.
					//
					re := regexp.MustCompile(`^#.*\r?\n`)
					jsonStr = re.ReplaceAllString(jsonStr, "")
					result := putRequest("http://localhost:9697/api/modal", strings.NewReader(jsonStr))
					fmt.Printf("%s", result[1:len(result)-1])
				} else {
					//
					// This is a raw html template that needs the data combined to finish it.
					//
					re := regexp.MustCompile(`\r?\n`)
					jsonStr = re.ReplaceAllString(jsonStr, " ")
					renderC := RenderDialogContents(jsonStr, data)
					result := putRequest("http://localhost:9697/api/dialog", strings.NewReader(renderC))
					fmt.Printf("%s", result[1:len(result)-1])
				}
			}
		}
	} else {
		//
		// Wrong information was given. Tell the user how to use the program.
		//
		// TODO: Needs better help information.
		//
		fmt.Printf("\n\nNot enough information!\nYou have to give the name of the dialog you want to show and the list of data to use in it.\nIf the only argument given is 'list', then a json list of dialogs is given.\nIf the only argument is 'build', than an interactive builder for a modal dialog will guide you through making one.")
	}
}
