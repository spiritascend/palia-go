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
	"github.com/buger/jsonparser"
)

type VersionData map[string]struct {
	BaseLineVer bool `json:"BaseLineVer"`
	Files       []struct {
		URL  string `json:"URL"`
		Hash string `json:"Hash"`
	} `json:"Files"`
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func GetConfig() (LocalConfig, error) {
	var Ret LocalConfig
	configfile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
		return LocalConfig{}, errors.New("failed to open config file")
	}
	defer configfile.Close()

	decoder := json.NewDecoder(configfile)
	decoder.Decode(&Ret)

	return Ret, nil
}

func GetLocalBuildVersion() (string, error) {
	Config, err := GetConfig()

	if err != nil {
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

func NeedUpdate() (bool, map[string]string, error) {

	LocalBuildVersion, err := GetLocalBuildVersion()
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	Config, err := GetConfig()

	if err != nil {
		return false, map[string]string{}, errors.New("failed to get local config")
	}

	resp, err := http.Get("https://update.palia.com/manifest/PatchManifest.json")
	if err != nil {
		log.Fatal("Error:", err)
		return false, map[string]string{}, errors.New("failed to get patch manifest")
	}
	defer resp.Body.Close()

	reqbody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return false, map[string]string{}, errors.New("error reading response body")
	}

	NeededFiles := make(map[string]string, 0)
	GetLatestExe := make(map[string]string, 0)

	jsonparser.ObjectEach(reqbody, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {

		baseLineVer, _ := jsonparser.GetBoolean(value, "BaseLineVer")

		jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			fileURL, _ := jsonparser.GetString(val, "URL")
			trimmedFileURL := strings.ReplaceAll(strings.ReplaceAll(fileURL, "https://update.palia.com/val/", ""), "v"+strings.ReplaceAll(string(key), ".", "")+"/", "")

			if string(key) != LocalBuildVersion && baseLineVer {
				NeededFiles[trimmedFileURL] = fileURL
			}

			if !baseLineVer {

				if strings.HasSuffix(trimmedFileURL, ".exe") {
					GetLatestExe[trimmedFileURL] = fileURL
				}

				if strings.HasSuffix(trimmedFileURL, ".pak") {
					if !fileExists(Config.Path + "\\Palia\\Content\\Paks\\" + trimmedFileURL) {
						NeededFiles[trimmedFileURL] = fileURL
					}
				}

			}
		}, "Files")
		return nil
	})

	if len(NeededFiles) == 0 {
		return false, NeededFiles, nil
	}
	NeededFiles["PaliaClient-Win64-Shipping.exe"] = GetLatestExe["PaliaClient-Win64-Shipping.exe"]
	fmt.Printf("Pending Install: %v\n", NeededFiles)
	return true, NeededFiles, nil
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

func HandleBaseGameDownload(pconfig *LocalConfig, url string, button *widget.Button) {

	// Create a temporary directory for downloading and unzipping
	tempDir, err := os.MkdirTemp("D:\\Palia", "palia_temp_")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		panic(1)
	}
	defer os.RemoveAll(tempDir)

	resp, err := http.Get(url)
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

}

func DownloadUpdate(filesneeded map[string]string, button *widget.Button) error {
	pconfig, err := GetConfig()

	if err != nil {
		return err
	}

	for key, value := range filesneeded {
		button.SetText("Updating")

		if strings.HasSuffix(key, ".zip") {
			HandleBaseGameDownload(&pconfig, value, button)
			continue
		}

		resp, err := http.Get(value)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer resp.Body.Close()

		totalSize := resp.ContentLength

		var binaryfilepath string

		if strings.HasSuffix(key, ".exe") {

			binaryfilepath = filepath.Join(pconfig.Path, "\\Palia\\Binaries\\Win64\\"+key)
		}
		if strings.HasSuffix(key, ".pak") {

			binaryfilepath = filepath.Join(pconfig.Path, "\\Palia\\Content\\Paks\\"+key)
		}

		file, err := os.Create(binaryfilepath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
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
			fmt.Println("Error while copying file content:", err)
			return err
		}
	}

	button.Enable()
	button.SetText("Launch Game")

	return nil
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
