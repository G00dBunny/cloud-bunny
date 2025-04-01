/*
┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ 				buba because its the ingrediants like types? get it :(?                                                           │
│                                                                                                                    │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
*/

package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"k8s.io/client-go/kubernetes"
)


type Model struct {

  //  NOTE - bubba use :)
  spinner         spinner.Model
  // nsList          list.Model
  // podList         list.Model
  logView         viewport.Model
  logContent      string
  state           int
  width           int
	height          int
  nsIndex         int
  podIndex        int
	selectedPod     string
	selectedNs      string
  //  NOTE - typical k8 use
  clientset       *kubernetes.Clientset
  namespaces       []string
  pods            []string
  // badPods         []string
	err             error
}








