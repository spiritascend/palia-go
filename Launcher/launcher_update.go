package launcher

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/widget"
)

type ManifestStructure struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	ManifestURL string `json:"manifest_url"`
	PatchMethod string `json:"patch_method"`
	Entry       string `json:"entry"`
}

type LocalConfig struct {
	Path string `json:"Path"`
}

type BuildVersion struct {
	MajorVersion         int    `json:"MajorVersion"`
	MinorVersion         int    `json:"MinorVersion"`
	PatchVersion         int    `json:"PatchVersion"`
	Changelist           int    `json:"Changelist"`
	CompatibleChangelist int    `json:"CompatibleChangelist"`
	IsLicenseeVersion    int    `json:"IsLicenseeVersion"`
	IsPromotedBuild      int    `json:"IsPromotedBuild"`
	BranchName           string `json:"BranchName"`
}

func GetConfig() (LocalConfig, bool) {
	var Ret LocalConfig
	configfile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
		return LocalConfig{}, true
	}
	defer configfile.Close()

	decoder := json.NewDecoder(configfile)
	decoder.Decode(&Ret)

	return Ret, false
}

func GetLocalBuildVersion() (string, error) {
	Config, hiccup := GetConfig()

	if hiccup {
		return "", errors.New("failed to get local config")
	}

	var BuildInfo BuildVersion

	buildversioninfo, err := os.Open(Config.Path + "\\Build.version")
	if err != nil {
		log.Fatal(err)
	}
	defer buildversioninfo.Close()

	decoder := json.NewDecoder(buildversioninfo)
	decoder.Decode(&BuildInfo)

	return strings.TrimPrefix(BuildInfo.BranchName, "++Valeria+Release_"), nil
}

func GetPaliaManifest() (ManifestStructure, error) {

	var Ret ManifestStructure

	resp, err := http.Get("https://update.palia.com/manifest/Palia.json")
	if err != nil {
		fmt.Println("Error:", err)
		return ManifestStructure{}, err
	}
	defer resp.Body.Close()

	reqbody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ManifestStructure{}, err
	}

	err = json.Unmarshal(reqbody, &Ret)
	if err != nil {
		fmt.Println("Error unmarshaling Manifest:", err)
		return ManifestStructure{}, err
	}

	return Ret, nil
}
func NeedToUpdate() (bool, error) {

	LocalBuildVersion, err := GetLocalBuildVersion()
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	Manifest, err := GetPaliaManifest()

	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	if LocalBuildVersion == Manifest.Version {
		return false, nil
	} else {
		return true, nil
	}
}

func unzipFile(zipFilePath, destinationFolder string) error {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destinationFolder, file.Name)

		if file.FileInfo().IsDir() {
			// Create directory if it doesn't exist
			err := os.MkdirAll(path, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		// Create the file
		writer, err := os.Create(path)
		if err != nil {
			return err
		}

		// Open and copy the file content
		src, err := file.Open()
		if err != nil {
			writer.Close()
			return err
		}

		_, err = io.Copy(writer, src)
		writer.Close()
		src.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleDownload(button *widget.Button) {
	button.SetText("Updating")
	pconfig, cerr := GetConfig()

	if cerr {
		fmt.Println(cerr)
		panic(1)
	}

	Manifest, err := GetPaliaManifest()
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	// Create a temporary directory for downloading and unzipping
	tempDir, err := os.MkdirTemp("D:\\Palia", "palia_temp_")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		panic(1)
	}
	defer os.RemoveAll(tempDir)

	resp, err := http.Get(Manifest.URL)
	if err != nil {
		fmt.Println("Error Getting Palia Game Zip File", err)
		panic(1)
	}
	defer resp.Body.Close()

	totalSize := resp.ContentLength

	zipFilePath := filepath.Join(tempDir, "game.zip")
	file, err := os.Create(zipFilePath)
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		panic(1)
	}
	defer file.Close()

	progress := &ProgressWriter{
		Writer:      io.MultiWriter(file),
		TotalSize:   totalSize,
		CurrentSize: 0,
		Button:      button,
	}

	_, err = io.Copy(progress, resp.Body)
	if err != nil {
		fmt.Println("Error while copying zip content:", err)
		panic(1)
	}

	destinationFolder := pconfig.Path

	button.SetText("Unzipping")
	err = unzipFile(zipFilePath, destinationFolder)
	if err != nil {
		fmt.Println("Error while unzipping:", err)
		panic(1)
	}

	button.Enable()
	button.SetText("Launch Game")
}

type ProgressWriter struct {
	Writer      io.Writer
	TotalSize   int64
	CurrentSize int64
	Button      *widget.Button
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.CurrentSize += int64(n)

	progress := float64(pw.CurrentSize) / float64(pw.TotalSize) * 100.0

	pw.Button.Text = fmt.Sprintf("Downloading: %.2f%%", progress)
	pw.Button.Refresh()

	return pw.Writer.Write(p)
}
