package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)


var (
	logViewStyle = 	lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#3C91E6")).
					Padding(1, 2)
)

const (
	stateNamespaceList = iota
	stateLoadingPods
	statePodList
	stateLoadingLogs
	statePodLogs
)


func StartUi(pods []string, namespaces []string, clientset *kubernetes.Clientset) error{

	m := initialModel(pods, namespaces, clientset)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()

	return err
}


func initialModel(pods []string, namespaces []string, clientset *kubernetes.Clientset) Model{
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))


	logView := viewport.New(0,0)
	logView.Style = logViewStyle

	return Model{
		spinner 	:	s,          
  		logView 	:   logView,      	
  		state 		:	stateNamespaceList,
  		clientset 	:   clientset,
  		namespaces 	:   namespaces,
  		pods  		:   pods,
		nsIndex		: 	0,
		podIndex	: 	0,
	}

}

func (m Model) Init() tea.Cmd {
	return nil
  }
  
func (m Model)  Update(msg tea.Msg) (tea.Model, tea.Cmd){
	return nil, nil
}

func (m Model)	View()	string {
	return ""
}



