package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbertschler/blocks/html"

	"github.com/Wes974/inventory/lib/guiapi"
	"github.com/Wes974/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["viewPart"] = viewPartAction
	guiapi.DefaultHandler.Functions["newPart"] = newPartAction
	guiapi.DefaultHandler.Functions["editPart"] = editPartAction
	guiapi.DefaultHandler.Functions["savePart"] = savePartAction
	guiapi.DefaultHandler.Functions["deletePart"] = deletePartAction
	guiapi.DefaultHandler.Functions["checkout"] = checkoutAction
}

func partPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/part/")
	part, err := parts.ByID(id)
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(viewPartBlock(part))
	err = html.Render(page, w)
	if err != nil {
		log.Println(err)
	}
}

func newPartAction(_ json.RawMessage) (*guiapi.Result, error) {
	return guiapi.Replace("#container", editPartBlock(nil, ""))
}

func editPartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", editPartBlock(part, ""))
}

func editPartBlock(p *parts.Part, code string) html.Block {
	isNew := false
	if p == nil {
		isNew = true
		p = &parts.Part{}
	}
	if code != "" {
		p.Code = code
	}
	cancelAction := fmt.Sprintf("guiapi('viewPart', '%s')", p.ID())
	saveAction := "sendForm('savePart', '.ga-edit-part')"
	return html.Div(nil,
		html.Div(nil,
			html.Button(html.Class("ui button").
				Attr("onclick", cancelAction),
				html.Text("Cancel"),
			),
			html.Button(html.Class("ui green button").
				Attr("onclick", saveAction),
				html.Text("Save"),
			),
		),
		html.Div(html.Class("ui form"),
			html.Input(html.Type("hidden").Name("New").Value(isNew).Class("ga-edit-part")),
			html.Input(html.Type("hidden").Name("ID").Value(p.ID()).Class("ga-edit-part")),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Code")),
				html.Input(html.Type("Text").Name("Code").Value(p.Code).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Location")),
				html.Input(html.Type("Text").Name("Location").Value(p.Location).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Family")),
				html.Input(html.Type("Text").Name("Family").Value(p.Family).Class("ga-edit-part")),
			),
		),
	)
}

func savePartAction(args json.RawMessage) (*guiapi.Result, error) {
	type input struct {
		ID       string
		New      string
		Code     string
		Location string
		Family   string
	}
	var in input
	err := json.Unmarshal(args, &in)
	if err != nil {
		return nil, err
	}
	var p *parts.Part
	if in.New == "true" {
		p, err = parts.Create()
	} else {
		p, err = parts.ByID(in.ID)
	}
	if err != nil {
		return nil, err
	}
	p.Code = in.Code
	p.Location = in.Location
	p.Family = in.Family

	// store part
	err = parts.Store(p)
	if err != nil {
		return nil, err
	}

	return guiapi.Replace("#container", viewPartBlock(p))
}

func deletePartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	err = parts.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Redirect("/")
}

func viewPartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(id)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return guiapi.Redirect("/")
	}
	return guiapi.Replace("#container", viewPartBlock(part))
}

func checkoutAction(args json.RawMessage) (*guiapi.Result, error) {
	var data struct {
		ID       string
	}
	err := json.Unmarshal(args, &data)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(data.ID)
	if err != nil {
		return nil, err
	}

	err = parts.Store(part)
	if err != nil {
		return nil, err
	}

	if part == nil {
		return guiapi.Redirect("/")
	}
	return guiapi.Replace("#container", viewPartBlock(part))
}

func viewPartBlock(p *parts.Part) html.Block {
	editAction := fmt.Sprintf("guiapi('editPart', '%s')", p.ID())
	deleteAction := fmt.Sprintf("guiapi('deletePart', '%s')", p.ID())

	var rows html.Blocks
	r := func(k, v string) html.Block {
		return html.Elem("tr", nil,
			html.Elem("td", nil, html.Text(k)),
			html.Elem("td", nil, html.Text(v)),
		)
	}
	rows.Add(r("Code", p.Code))
	rows.Add(r("Location", p.Location))
	rows.Add(r("Family", p.Family))

	return html.Div(nil,
		html.Div(nil,
			html.A(html.Href("/"),
				html.Button(html.Class("ui button"),
					html.Text("< List"),
				),
			),
			html.Button(html.Class("ui button").
				Attr("onclick", editAction),
				html.Text("Edit"),
			),
			html.Button(html.Class("ui red button").
				Attr("onclick", deleteAction),
				html.Text("Delete"),
			),
		),
		html.H1(nil, html.Text(p.Code)),
		html.Div(html.Class("ui form"),
			html.Input(html.Type("hidden").Name("ID").Value(p.ID()).Class("ga-checkout")),
			// TODO: Change checkout ?
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Inventory")),
				// html.Input(html.Type("Text").Value(p.Quantity).Attr("disabled", true)),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Checkout")),
				// html.Input(html.Type("Text").Name("Quantity").Value("0").Class("ga-checkout")),
			),
			html.Button(html.Class("ui yellow button").
				Attr("onclick", "sendForm('checkout', '.ga-checkout')"),
				html.Text("Checkout Parts"),
			),
		),
		html.Elem("table", html.Class("ui celled table"),
			html.Elem("tbody", nil,
				rows,
			),
		),
	)
}
