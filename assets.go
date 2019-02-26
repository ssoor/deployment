package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func buildAssets(assetsSplit []Assets) error {
	for _, assets := range assetsSplit {
		// copy file
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if match, err := filepath.Match(assets.Source, path); nil != err {
				return err
			} else if !match {
				return nil
			}

			baseDir := path
			matchDir := assets.Source

			for {
				baseDir = filepath.Dir(baseDir)
				matchDir = filepath.Dir(matchDir)

				if matchDir == baseDir || "." == baseDir {
					break
				}
			}

			targetPath, err := filepath.Rel(baseDir, path)
			if nil != err {
				return err
			}
			targetPath = assets.Target + "/" + targetPath

			log.Print("COPY: ", path, " => ", targetPath)
			// log.Print("Match: [", match, "] ", path, " => ", assets.Source)

			if srcinfo, err := os.Stat(path); err != nil {
				return err
			} else if err = os.MkdirAll(filepath.Dir(targetPath), srcinfo.Mode()); err != nil {
				return err
			}

			if info.IsDir() {
				return CopyDir(path, targetPath)
			}

			return CopyFile(path, targetPath)
		})

		if nil != err {
			return err
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err := os.Mkdir(dst, srcinfo.Mode()); nil != err {
		if !os.IsExist(err) {
			return err
		}
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}

	return nil
}
