package desktop

import (
	"errors"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var a fyne.App
var configWindow fyne.Window

func Init(icon []byte) fyne.App {
	a = app.New()
	setupConfigWindow()

	a.SetIcon(fyne.NewStaticResource("Icon.png", icon))

	if desk, ok := a.(desktop.App); ok {
		showItem := fyne.NewMenuItem("Config", func() {
			configWindow.Show()
		})
		quitItem := fyne.NewMenuItem("Quit", func() {
			a.Quit()
		})
		m := fyne.NewMenu("MyApp", showItem, quitItem)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(a.Icon())
	}

	a.SendNotification(fyne.NewNotification("lylink-jellyfin", "LyLink is running in the background"))

	return a
}

func setupConfigWindow() {
	configWindow = a.NewWindow("Config")

	container := container.New(
		layout.NewFormLayout(),
	)
	container.Add(widget.NewLabel("Service port"))
	portE := widget.NewEntry()
	port := "8080"
	portE.SetText(port)
	portE.Validator = func(s string) error {
		if s == "" {
			return nil
		} else {
			port, err := strconv.Atoi(s)
			if port < 0 || port > 65535 {
				err = errors.New("port must be between 0 and 65535")
			}
			return err
		}

	}
	container.Add(portE)

	container.Add(widget.NewLabel("Jellyfin server url"))

	serverUrlE := widget.NewEntry()
	serverUrl := "http://localhost:8096"
	serverUrlE.Text = serverUrl
	serverUrlE.Validator = func(s string) error {
		if s == "" {
			return nil
		} else {
			uri, err := url.ParseRequestURI(s)
			if uri.Host == "" {
				err = errors.New("invalid url")
			}
			return err
		}

	}
	container.Add(serverUrlE)

	configWindow.SetContent(container)

	configWindow.Show()
}
