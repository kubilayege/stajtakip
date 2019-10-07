package templates

import (
  "html/template"
)

var tpl = template.Must(template.ParseFiles("templates/main.html"))

type Main struct {
  Title string
  Alert string
  AlertMessage string
  Vars interface{}
}

func (m Main) Error(msg string) Main {
  m.Alert = "danger"
  m.AlertMessage = msg
  return m
}

func (m Main) Warning(msg string) Main {
  m.Alert = "warning"
  m.AlertMessage = msg
  return m
}

func (m Main) Info(msg string) Main {
  m.Alert = "info"
  m.AlertMessage = msg
  return m
}

func NewMain(title string) Main {
  return Main{
    Title: title,
    Alert: "",
    AlertMessage: "",
  }
}

func Load(filenames ...string) *template.Template {
  return template.Must(template.Must(tpl.Clone()).ParseFiles(filenames...))
}
