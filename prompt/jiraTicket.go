package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/G00dBunny/cloud-bunny/jiraBed"
)

func PromptUserForTicketConfirmation(details *jiraBed.TicketDetail) bool {
	fmt.Println("\n=== JIRA TICKET ===")
	fmt.Println("Title:", details.Summary)
	fmt.Println("\nDescription:")
	fmt.Println(details.Description)
	fmt.Println("\nIssue Type:", details.IssueType)		//Hard coded for right now
	fmt.Println("Priority:", details.Priority)			//Hard coded for right now
	
	fmt.Println("\n######Need to edit???########### (y/n):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	
	if strings.TrimSpace(strings.ToLower(input)) == "y" {

		fmt.Println("\nCurrent Summary:", details.Summary)
		fmt.Println("Press Enter to keep current or Title")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			details.Summary = input
		}
		
		fmt.Println("\nCurrent Description:")
		fmt.Println(details.Description)
		fmt.Println("\nPress Enter to keep current or type 'edit' to modify:")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "edit" {
			fmt.Println("Enter new description (type 'END' at end to switch lines):")
			var lines []string
			for {
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == "END" {
					break
				}
				lines = append(lines, line)
			}
			if len(lines) > 0 {
				details.Description = strings.Join(lines, "\n")
			}
		}
		
		// NOTE : hardcoded should be changed
		fmt.Println("\nCurrent Issue Type:", details.IssueType)
		fmt.Println("Press Enter to keep current or type new issue type (Bug, Task, Story):")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			details.IssueType = input
		}
		
		// NOTE : hardcoded should also be changed
		fmt.Println("\nCurrent Priority:", details.Priority)
		fmt.Println("Press Enter to keep current or type new priority (Highest, High, Medium, Low, Lowest):")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			details.Priority = input
		}
	}
	
	fmt.Println("\nCreate this ticket? (y/n):")
	input, _ = reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}