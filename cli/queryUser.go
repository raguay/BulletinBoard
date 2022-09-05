package main

import (
	"net/http"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"fmt"
	"os"
	"github.com/aymerick/raymond"
	"path/filepath"
	"path"
	"regexp"
	"encoding/json"
	"github.com/charmbracelet/bubbletea"
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
// Struct:		model
//
// description: The structure for the bubbletea interface for building a dialog.
//
type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    selected map[int]struct{}   // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our list of acctions
		choices:  []string{"Add Item", "Add Button", "Save", "Quit"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

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

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m model) View() string {
    // The header
    s := "What do you want to do?\n\n"

    // Iterate over our choices
    for i, choice := range m.choices {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x" // selected!
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
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
func main()  {
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
		progHome,_ := os.Executable()
		progHome = filepath.Dir(progHome)
		templates1 := filepath.Join(progHome,"../Resources/dialogs")
		templates2 := filepath.Join(home,".config/scriptserver/dialogs")
		
		//
		// Process the command or the template.
		//
		switch dialog {
			case "list": {
				//
				// Give the user a json list of dialogs in the program
				// area and in the user directory.
				//
				var nlist []string
				file, _ := os.Open(templates1)
				dlist,_ := file.Readdirnames(0) // 0 to read all files and folders
   				for _, name := range dlist {
       				nlist = append(nlist,name)
   				}
				file.Close()
				file, _ = os.Open(templates2)
				dlist,_ = file.Readdirnames(0)
					for _, name := range dlist {
					nlist = append(nlist,name)
				}
				file.Close()
				nlist = Map(nlist, FilenameWithoutExtension)
				pjson, err := json.Marshal(nlist)
   				if err != nil {
       				log.Fatal("Cannot encode to JSON ", err)
   				}
				fmt.Printf("{ \"dialogs\": %s}\n", pjson)
			}
			case "build": {
				//
				// We are going to build a dialog.
				//
				p := tea.NewProgram(initialModel())
    			if err := p.Start(); err != nil {
        			fmt.Printf("Alas, there's been an error: %v", err)
        			os.Exit(1)
    			}
			}
			default: {
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
				templatefile1 := filepath.Join(templates1, fmt.Sprintf("%s.json",dialog))
				templatefile2 := filepath.Join(templates2, fmt.Sprintf("%s.json",dialog))
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
					jsonStr = re.ReplaceAllString(jsonStr,"")
					result := putRequest("http://localhost:9697/api/modal", strings.NewReader(jsonStr))
					fmt.Printf("%s",result[1:len(result)-1])
				} else {
					//
					// This is a raw html template that needs the data combined to finish it.
					//
					re := regexp.MustCompile(`\r?\n`)
					jsonStr = re.ReplaceAllString(jsonStr, " ")
					renderC := RenderDialogContents(jsonStr, data)
					result := putRequest("http://localhost:9697/api/dialog", strings.NewReader(renderC))
					fmt.Printf("%s",result[1:len(result)-1])
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
