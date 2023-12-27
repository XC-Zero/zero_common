package zip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"github.com/pkg/errors"
	"log"
	"strings"

	"io"
	"os"
	"path"
)

func UnzipTarGZ(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return errors.WithStack(err)
	}
	dirPath := path.Dir(filepath)

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(dirPath, header.Name), 0755); err != nil {
				return errors.WithStack(err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(path.Join(dirPath, header.Name))
			if err != nil {
				return errors.WithStack(err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.WithStack(err)
			}
			outFile.Close()
		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}

	}
	return nil
}

func UnzipGZ(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return errors.WithStack(err)
	}
	p, _ := strings.CutSuffix(filepath, ".gz")
	output, err := os.Create(p)
	if err != nil {
		return errors.WithStack(err)
	}
	// 将 tar 文件内容拷贝到输出文件
	_, err = io.Copy(output, reader)
	if err != nil {
		return errors.WithStack(err)
	}
	output.Close()
	return nil
}

func UnzipZIP(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return errors.WithStack(err)
	}
	for i := range reader.File {

		open, err := reader.File[i].Open()
		if err != nil {
			return err
		}
		output, err := os.Create(path.Join(path.Dir(filepath), reader.File[i].Name))
		if err != nil {
			return errors.WithStack(err)
		}
		// 将 tar 文件内容拷贝到输出文件
		_, err = io.Copy(output, open)
		if err != nil {
			return errors.WithStack(err)
		}
		output.Close()

	}
	return nil
}
