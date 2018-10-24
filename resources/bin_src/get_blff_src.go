package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	blffDownloadURL         = "https://github.com/syncore/blff/archive/master.zip"
	blffDownloadFilename    = "blff_src.zip"
	nugetDownloadURL        = "https://dist.nuget.org/win-x86-commandline/latest/nuget.exe"
	nugetDownloadFilename   = "nuget.exe"
	extractDistPkgFilesFlag = "extractDistPkgs"
	distPkgDir              = "blffDistPkg"
	binDir                  = "..\\bin"
)

var unzipDistPkgFilesArg bool

func init() {
	flag.BoolVar(&unzipDistPkgFilesArg, extractDistPkgFilesFlag, false, "specify flag to extract blff dist packages containing blff binaries")
}

func main() {
	flag.Parse()
	if unzipDistPkgFilesArg {
		if err := extractBlffDistPackages(); err != nil {
			fmt.Printf("Unable to extract blff dist packages: %v", err)
			return
		}
		fmt.Printf("Successfully extracted blff distpkgs to: %s\n", filepath.Join(getExecutingPath(), binDir))
		return
	}
	if err := getBlffSource(); err != nil {
		fmt.Printf("Failed to get blff source code: %v", err)
		return
	}
	fmt.Println("Successfully downloaded blff source code and supporting utilities.")
}

func getBlffSource() error {
	execDir := getExecutingPath()
	if err := download(filepath.Join(execDir, blffDownloadFilename), blffDownloadURL); err != nil {
		return fmt.Errorf("Unable to download blff source: %v\n", err)
	}
	if err := unzip(blffDownloadFilename, execDir); err != nil {
		return fmt.Errorf("Unable to unzip downloaded blff source: %v\n", err)
	}
	if err := os.Rename(filepath.Join(execDir, "blff-master"), filepath.Join(execDir, "blff_src")); err != nil {
		return fmt.Errorf("Unable to rename blff-master directory to blff_src: %v\n", err)
	}
	if err := download(filepath.Join(execDir, nugetDownloadFilename), nugetDownloadURL); err != nil {
		return fmt.Errorf("Unable to download latest nuget: %v\n", err)
	}
	return nil
}

func getExecutingPath() string {
	executingPath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Fatal error: unable to determine execution path: %v", err))
	}
	return filepath.Dir(executingPath)
}

func extractBlffDistPackages() error {
	execDir := getExecutingPath()
	dpkgDir := filepath.Join(execDir, distPkgDir)
	var distPkgFiles []string
	if err := filepath.Walk(dpkgDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			distPkgFiles = append(distPkgFiles, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("Unable to get blff distpkg directory listing: %v\n", err)
	}
	for _, dpf := range distPkgFiles {
		if err := unzip(dpf, filepath.Join(execDir, binDir)); err != nil {
			return fmt.Errorf("Unable to extract blff distpkg: %s: %v\n", dpf, err)
		}
		if err := os.Remove(dpf); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("Unable to remove blff distpkg: %s: %v\n", dpf, err)
		}
	}
	if err := os.Remove(dpkgDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unable to remove blff distpkg directory: %s: %v\n", dpkgDir, err)
	}
	return nil
}

func download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func unzip(source, destination string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, zf := range r.File {
		if err := unzipFile(zf, destination); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(zf *zip.File, destination string) error {
	fpath := filepath.Join(destination, zf.Name)
	if !strings.HasPrefix(fpath, filepath.Clean(destination)) {
		return fmt.Errorf("%s: illegal file path", destination)
	}
	if strings.HasSuffix(zf.Name, "/") {
		if err := os.MkdirAll(fpath, 0755); err != nil {
			return fmt.Errorf("%s: making directory: %v", fpath, err)
		}
		return nil
	}
	rc, err := zf.Open()
	if err != nil {
		return fmt.Errorf("%s: open compressed file: %v", zf.Name, err)
	}
	defer rc.Close()
	if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		return fmt.Errorf("%s: making directory for file: %v", fpath, err)
	}
	out, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("%s: creating new file: %v", fpath, err)
	}
	defer out.Close()
	if err = out.Chmod(zf.FileInfo().Mode()); err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("%s: changing file mode: %v", fpath, err)
	}
	_, err = io.Copy(out, rc)
	if err != nil {
		return fmt.Errorf("%s: writing file: %v", fpath, err)
	}
	return nil
}
