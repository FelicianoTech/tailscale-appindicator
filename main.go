package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
)

// checks connectivity using shell's ping.
// Can move to a Go lib functionality in the future, which will require sudo
// https://github.com/go-ping/ping
// rights but this will be running as root anyway
func connectionWorks() bool {

	cmd := exec.Command("ping", "-c", "3", "100.101.102.103")
	cmd.Start()
	err := cmd.Wait()

	if err != nil {
		return false
	}

	return true
}

// For now this function returns a single IP. In the future it will return a status struct.
func pullStatus() string {

	statusOutput, err := exec.Command("tailscale", "status", "--json").Output()
	if err != nil {
		log.Fatal(err)
	}

	var statusJSON map[string]interface{}
	var ip string

	err = json.Unmarshal(statusOutput, &statusJSON)
	if err != nil {
		log.Fatal(err)
	}

	ip = statusJSON["TailscaleIPs"].([]interface{})[0].(string)

	return ip
}

func main() {

	gtk.Init(nil)

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}

	indicator := appindicator.NewWithPath("indicator-tailscale", "tailscale-icon", appindicator.CategoryApplicationStatus, "/home/felicianotech/Repos/felicianotech/tailscale-appindicator")
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)
	indicator.SetLabel(pullStatus(), "")

	itemStatus, err := gtk.MenuItemNewWithLabel("Status: unknown")
	if err != nil {
		log.Fatal(err)
	}

	sep1, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		log.Fatal(err)
	}

	itemAdmin, err := gtk.MenuItemNewWithLabel("Open Tailscale Admin")
	if err != nil {
		log.Fatal(err)
	}
	_ = itemAdmin.Connect("activate", func() {
		exec.Command("xdg-open", "https://login.tailscale.com/admin/machines").Start()
	})
	if err != nil {
		log.Fatal(err)
	}

	itemDocs, err := gtk.MenuItemNewWithLabel("Open Tailscale Docs")
	if err != nil {
		log.Fatal(err)
	}
	_ = itemDocs.Connect("activate", func() {
		exec.Command("xdg-open", "https://tailscale.com/kb/").Start()
	})
	if err != nil {
		log.Fatal(err)
	}

	sep2, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		log.Fatal(err)
	}

	itemExit, err := gtk.MenuItemNewWithLabel("Exit AppIndicator")
	if err != nil {
		log.Fatal(err)
	}
	_ = itemExit.Connect("activate", func() {
		gtk.MainQuit()
	})
	if err != nil {
		log.Fatal(err)
	}

	menu.Add(itemStatus)
	menu.Add(sep1)
	menu.Add(itemAdmin)
	menu.Add(itemDocs)
	menu.Add(sep2)
	menu.Add(itemExit)
	menu.ShowAll()

	go func() {

		for {
			<-time.After(time.Second * 5)

			if connectionWorks() {
				itemStatus.SetLabel("Status: connected")
			} else {
				itemStatus.SetLabel("Status: disconnected")
			}
		}
	}()

	gtk.Main()
}
