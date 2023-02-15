package main

import (
	"archive/zip"
	"io"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Installer")

	// create a text label
	label := widget.NewLabel("Click the button below to install the program")

	// create a button that will extract the contents of the zip package to a specified directory
	button := widget.NewButton("Install", func() {
		zipReader, err := zip.OpenReader("program.zip")
		if err != nil {
			label.SetText("Error opening zip file")
			return
		}
		defer zipReader.Close()

		// specify the directory to extract the files to
		extractDir := "./program"

		for _, file := range zipReader.File {
			if file.FileInfo().IsDir() {
				os.MkdirAll(extractDir+"/"+file.Name, file.Mode())
				continue
			}

			outFile, err := os.OpenFile(extractDir+"/"+file.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				label.SetText("Error extracting file: " + file.Name)
				return
			}
			defer outFile.Close()

			fileReader, err := file.Open()
			if err != nil {
				label.SetText("Error opening file: " + file.Name)
				return
			}
			defer fileReader.Close()

			if _, err = io.Copy(outFile, fileReader); err != nil {
				label.SetText("Error extracting file: " + file.Name)
				return
			}
		}

		label.SetText("Installation successful!")
	})

	// create a container to hold the label and button
	content := container.NewVBox(
		label,
		button,
	)

	// set the window content
	w.SetContent(content)

	// show the window
	w.ShowAndRun()
}
