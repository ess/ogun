// Package fs provides services for dealing with file systems
package fs

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/ess/ogun/pkg/ogun"
)

var Root = afero.NewOsFs()

var CreateDir = func(path string, mode os.FileMode) error {
	if !FileExists(path) {
		err := Root.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

var FileExists = func(path string) bool {
	_, err := Root.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

var DirectoryExists = func(path string) bool {
	if !FileExists(path) {
		return false
	}

	if !IsDir(path) {
		return false
	}

	return true
}

var IsDir = func(path string) bool {
	info, err := Root.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func Walk(path string, walkFunc filepath.WalkFunc) error {
	return afero.Walk(Root, path, walkFunc)
}

func Copy(path string, targetPath string, mode os.FileMode) error {
	infile, err := Root.Open(path)
	if err != nil {
		return err
	}
	defer infile.Close()

	if FileExists(targetPath) {
		rmerr := Root.Remove(targetPath)
		if rmerr != nil {
			return rmerr
		}
	}

	outfile, err := Root.Create(targetPath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := outfile.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(outfile, infile)
	if err != nil {
		return err
	}

	err = outfile.Sync()
	return err
}

func DirCopy(source string, target string) error {
	if !FileExists(source) {
		return fmt.Errorf("%s does not exist", source)
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		targetPath := target + strings.TrimPrefix(path, source)

		if info.IsDir() {
			return CreateDir(targetPath, info.Mode())
		} else {
			return Copy(path, targetPath, info.Mode())
		}
	}

	return Walk(source, walkFunc)
}

var Tar = func(src string, destination string) error {

	destfile, err := Root.Create(destination)
	if err != nil {
		return err
	}

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return err
	}

	gzw := gzip.NewWriter(destfile)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}

func CreateBuildLog(application string, name string) (afero.File, error) {
	logName := "build-" + name + ".log"
	logPath := applicationPath(ogun.Application{Name: application}) + "/shared/build_logs/"

	err := CreateDir(logPath, 0755)
	if err != nil {
		return nil, err
	}

	return Root.Create(logPath + logName)
}

func ReadDir(path string) ([]os.FileInfo, error) {
	return afero.ReadDir(Root, path)
}

func Basename(path string) string {
	return filepath.Base(path)
}

func Stat(path string) (os.FileInfo, error) {
	return Root.Stat(path)
}

func Executable(path string) bool {
	info, _ := Stat(path)

	return (info.Mode()&0100 == 0100)
}
