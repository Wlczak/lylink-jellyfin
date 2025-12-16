package desktop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Wlczak/lylink-jellyfin/api"
	"github.com/Wlczak/lylink-jellyfin/config"
	"github.com/Wlczak/lylink-jellyfin/logs"
	"github.com/gin-gonic/gin"
)

var a fyne.App
var configWindow fyne.Window
var router *gin.Engine
var server *http.Server
var aboutWindow fyne.Window
var updateButton *widget.Button

func Init(icon []byte, r *gin.Engine, srv *http.Server) fyne.App {
	router = r
	server = srv

	a = app.New()
	setupConfigWindow()

	a.SetIcon(fyne.NewStaticResource("Icon.png", icon))

	if desk, ok := a.(desktop.App); ok {
		showItem := fyne.NewMenuItem("Config", func() {
			configWindow.Show()
		})
		aboutItem := fyne.NewMenuItem("About", func() {
			aboutWindow = a.NewWindow("About")
			l := container.NewVBox()
			versionLabel := widget.NewLabel("lylink-jellyfin v0.0.1")

			updateButton = widget.NewButton("Check for updates", func() {
				updateButton.SetText("Checking...")
				updateAvailable, newVersion, _ := hasUpdate()
				if updateAvailable {
					a.SendNotification(&fyne.Notification{Title: "New version " + newVersion + " available"})
				}
				updateButton.SetText("Check for updates")
			})

			l.Add(versionLabel)
			l.Add(updateButton)
			aboutWindow.SetContent(l)
			aboutWindow.Show()
		})
		quitItem := fyne.NewMenuItem("Quit", func() {
			a.Quit()
		})
		m := fyne.NewMenu("MyApp", showItem, aboutItem, quitItem)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(a.Icon())
	}

	a.SendNotification(fyne.NewNotification("lylink-jellyfin", "LyLink is running in the background"))

	updateAvailable, versionName, _ := hasUpdate()
	if updateAvailable {
		a.SendNotification(&fyne.Notification{Title: "New version " + versionName + " available"})
	}

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
			d := dialog.NewError(errors.New("config is invalid"), configWindow)
			d.MinSize()
			d.Show()
			return
		}
		serverUrl = serverUrlE.Text
		serverUrl = strings.TrimSuffix(serverUrl, "/")
		port, err := strconv.Atoi(portE.Text)
		if err != nil {
			port = 0
			zap.Error(err.Error())
		}
		config.Port = port
		config.JellyfinServerUrl = serverUrl
		err = config.Save()

		if err != nil {
			zap.Error(err.Error())
			d := dialog.NewError(err, configWindow)
			d.MinSize()
			d.Show()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = server.Shutdown(ctx)

		if err != nil {
			zap.Error(err.Error())
		}

		cancel()

		server = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}

		go api.RunHttpServer(server)

		configWindow.Hide()
	})

	vbox := container.NewVBox(
		form, layout.NewSpacer(), submit)

	configWindow.SetContent(vbox)

	size := vbox.MinSize()
	size.Width = size.Width + 200
	size.Height = size.Height + 100
	configWindow.Resize(size)
}

func hasUpdate() (bool, string, error) {
	resp, err := http.Get("https://api.github.com/repos/wlczak/lylink-jellyfin/releases/latest")
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return false, "", err
	}
	return release.TagName != "v0.0.1", release.TagName, nil
}
