package buildcontext

import (
	"archive/tar"
	"io"
	"os"

	"github.com/pkg/errors"
)

const Dockerfile = "Dockerfile"

type BuildContext struct {
	out       io.Writer
	tarWriter *tar.Writer
}

func NewBuildContext(out io.Writer) *BuildContext {
	tarWriter := tar.NewWriter(out)

	return &BuildContext{
		out:       out,
		tarWriter: tarWriter,
	}
}

func (bc *BuildContext) WriteDockerfile(data []byte) error {
	if err := bc.tarWriter.WriteHeader(&tar.Header{
		Name:     Dockerfile,
		Size:     int64(len(data)),
		Typeflag: tar.TypeReg,
	}); err != nil {
		return errors.Wrap(err, "Dockerfile: tar header write failed")
	}

	_, err := bc.tarWriter.Write(data)
	if err != nil {
		return errors.Wrapf(err, "Dockerfile: tar write failed")
	}
	return nil
}

func (bc *BuildContext) WriteFileData(data []byte, path string) error {
	if err := bc.tarWriter.WriteHeader(&tar.Header{
		Name:     path,
		Size:     int64(len(data)),
		Typeflag: tar.TypeReg,
	}); err != nil {
		return errors.Wrapf(err, "%s: tar header write failed", path)
	}

	_, err := bc.tarWriter.Write(data)
	if err != nil {
		return errors.Wrapf(err, "%s: tar write failed", path)
	}
	return nil
}

func (bc *BuildContext) WriteFile(fileInfo os.FileInfo, reader io.Reader, path string) error {
	hdr, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return errors.Wrapf(err, "%s: building tar header failed", path)
	}
	hdr.Name = path

	if err := bc.tarWriter.WriteHeader(hdr); err != nil {
		return errors.Wrapf(err, "%s: tar header write failed", path)
	}

	_, err = io.Copy(bc.tarWriter, reader)
	if err != nil {
		return errors.Wrapf(err, "%s: copy failed", path)
	}

	return nil
}

func (bc *BuildContext) Close() error {
	return bc.tarWriter.Close()
}
