package desktop

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wlczak/lylink-jellyfin/config"
	"github.com/Wlczak/lylink-jellyfin/logs"
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
	zap := logs.GetLogger()
	config := config.GetConfig()

	configWindow = a.NewWindow("Config")
	configWindow.SetCloseIntercept(func() {
		configWindow.Hide()
	})

	form := container.New(
		layout.NewFormLayout(),
	)
	form.Add(widget.NewLabel("Service port"))
	portE := widget.NewEntry()
	port := config.Port
	portE.SetText(fmt.Sprintf("%d", port))
	portE.Validator = func(s string) error {
		port, err := strconv.Atoi(s)
		if port < 1 || port > 65535 {
			err = errors.New("port must be between 0 and 65535")
		}
		return err
	}
	form.Add(portE)

	form.Add(widget.NewLabel("Jellyfin server url"))

	serverUrlE := widget.NewEntry()
	serverUrl := config.JellyfinServerUrl
	serverUrlE.Text = serverUrl
	serverUrlE.Validator = func(s string) error {
		uri, err := url.ParseRequestURI(s)
		if err != nil {
			return errors.New("invalid url")
		}
		if uri.Host == "" {
			err = errors.New("invalid url")
		}
		return err
	}
	form.Add(serverUrlE)

	submit := widget.NewButton("Submit", func() {
		if portE.Validate() != nil || serverUrlE.Validate() != nil {
			d := dialog.NewError(errors.New("Config is invalid"), configWindow)
			d.MinSize()
			d.Show()
			return
		}
		serverUrl = serverUrlE.Text
		port, err := strconv.Atoi(portE.Text)
		if err != nil {
			port = 0
			zap.Error(err.Error())
		}
		config.Port = port
		config.JellyfinServerUrl = serverUrl
		config.Save()
		configWindow.Hide()
	})

	vbox := container.NewVBox(
		form, layout.NewSpacer(), submit)

	configWindow.SetContent(vbox)

	size := vbox.MinSize()
	size.Width = size.Width + 200
	size.Height = size.Height + 100
	configWindow.Resize(size)

	configWindow.Show()
}
