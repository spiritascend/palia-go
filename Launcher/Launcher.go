package launcher

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func LaunchGame() bool {
	config, err := GetConfig()

	if err {
		fmt.Println(err)
		panic(1)
	}
	cmd := exec.Command(config.Path+"\\Palia\\Binaries\\Win64\\PaliaClient-Win64-Shipping.exe", "-console", "-log")

	cmd.Run()

	return true
}

func IntiializeLauncher() {
	a := app.New()
	w := a.NewWindow("Palia Launcher")

	progressbar := widget.NewProgressBar()

	var LaunchOrDownloadButton *widget.Button

	LaunchOrDownloadButton = widget.NewButton("{Button}", func() {
		NeedToUpdate, err := NeedToUpdate()

		if err != nil {
			fmt.Println(err)
			panic(1)
		}

		if !NeedToUpdate {
			go LaunchGame()
			LaunchOrDownloadButton.Text = "Launch Game"
			return
		} else {
			LaunchOrDownloadButton.Disable()
			LaunchOrDownloadButton.Text = "Update Game"
			go HandleDownload(LaunchOrDownloadButton)
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
