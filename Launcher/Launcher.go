package launcher

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func LaunchGame() bool {
	config, err := GetConfig()

	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	lperr := LaunchProcess(config.Path + "\\Palia\\Binaries\\Win64\\PaliaClient-Win64-Shipping.exe -console -log")

	return lperr != nil
}

func IntiializeLauncher() {
	a := app.New()
	w := a.NewWindow("Palia Launcher")

	progressbar := widget.NewProgressBar()

	var LaunchOrDownloadButton *widget.Button

	LaunchOrDownloadButton = widget.NewButton("{Button}", func() {
		needtoupdate, filesneeded, err := NeedUpdate()

		if !needtoupdate || err != nil {
			go LaunchGame()
			LaunchOrDownloadButton.Text = "Launch Game"
			return
		} else {

			downloaderrorChannel := make(chan error)

			LaunchOrDownloadButton.Disable()
			LaunchOrDownloadButton.Text = "Update Game"

			go func() {
				err := DownloadUpdate(filesneeded, LaunchOrDownloadButton)
				if err != nil {
					downloaderrorChannel <- err
					return
				}
				downloaderrorChannel <- nil
			}()

			err := <-downloaderrorChannel

			if err != nil {
				LaunchOrDownloadButton.Enable()
				LaunchOrDownloadButton.Text = "Launch Game"
			}

		}

	})

	LaunchOrDownloadButton.Text = "Launch Game"

	progressbar.Hide()
	content := container.NewVBox(progressbar, LaunchOrDownloadButton)
	BottumContainer := container.NewBorder(nil, content, nil, nil)

	tab2Content := container.NewVBox()

	HomeTab := container.NewTabItem("Home", BottumContainer)
	PrivateServerTab := container.NewTabItem("Private Server", tab2Content)

	tabs := container.NewAppTabs(HomeTab, PrivateServerTab)

	w.SetContent(tabs)

	w.Resize(fyne.NewSize(1000, 750))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
