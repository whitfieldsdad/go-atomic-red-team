package atomic_red_team

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/pkg/errors"
)

func CreateTarball(archivePath string, filePaths []string) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			errors.Wrapf(err, "failed to add file %s to tar writer", filePath)
		}
	}
	return nil
}

func addFileToTarWriter(path string, tarWriter *tar.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to stat file")
	}
	header := &tar.Header{
		Name:    path,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.Wrap(err, "failed to write tar header")
	}
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.Wrap(err, "failed to copy file to tar writer")
	}
	return nil
}
