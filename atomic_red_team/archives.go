package atomic_red_team

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"filippo.io/age"
	"github.com/pkg/errors"
)

func CreateTarball(writer io.Writer, filePaths []string) error {
	tarWriter, err := createTarball(writer, filePaths)
	if err != nil {
		return err
	}
	defer tarWriter.Close()
	return nil
}

func CreateTarballFile(archivePath string, filePaths []string) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()
	return CreateTarball(file, filePaths)
}

func CreateEncryptedTarball(writer io.Writer, filePaths []string, password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	tarWriter, err := createTarball(writer, filePaths)
	if err != nil {
		return err
	}
	defer tarWriter.Close()

	identity, err := age.NewScryptRecipient(password)
	if err != nil {
		return errors.Wrap(err, "failed to create scrypt identity")
	}
	w, err := age.Encrypt(tarWriter, identity)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt tarball")
	}
	defer w.Close()
	return nil
}

func CreateEncryptedTarballFile(archivePath string, filePaths []string, password string) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()
	return CreateEncryptedTarball(file, filePaths, password)
}

func createTarball(writer io.Writer, filePaths []string) (*tar.Writer, error) {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			errors.Wrapf(err, "failed to add file %s to tar writer", filePath)
		}
	}
	return tarWriter, nil
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

func ReadTarball(reader io.Reader) (*tar.Reader, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzipReader)
	return tarReader, nil
}

func ReadEncryptedTarball(reader io.Reader, password string) (*tar.Reader, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	id, err := age.NewScryptIdentity(password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create identity")
	}
	decryptedReader, err := age.Decrypt(reader, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt tarball")
	}
	return ReadTarball(decryptedReader)
}
