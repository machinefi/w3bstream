package resource

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/types"
)

var reserve = int64(100 * 1024 * 1024)

func Upload(ctx context.Context, f *multipart.FileHeader, id string) (root, fullName, sum string, err error) {
	l := types.MustLoggerFromContext(ctx)
	conf := types.MustUploadConfigFromContext(ctx)
	var (
		fr       io.ReadSeekCloser
		fw       io.WriteCloser
		filesize = int64(0)
	)

	_, l = l.Start(ctx, "Upload")
	defer l.End()

	root = filepath.Join(conf.Root, id)
	fullName = filepath.Join(conf.Root, id, f.Filename)

	if !IsDirExists(root) {
		if err = os.MkdirAll(root, 0777); err != nil {
			l.Error(err)
			return
		}
	}

	if fr, err = f.Open(); err != nil {
		return
	}
	defer fr.Close()

	if filesize, err = fr.Seek(0, io.SeekEnd); err != nil {
		l.Error(err)
		return
	}
	if filesize > conf.FileSizeLimit {
		err = errors.Wrap(err, "filesize over limit")
		l.Error(err)
		return
	}

	stat, err := disk.Usage(root)
	if stat == nil || stat.Free < uint64(filesize+reserve) {
		err = errors.Wrap(err, "disk limited")
		l.Error(err)
		return
	}
	_, err = fr.Seek(0, io.SeekStart)
	if err != nil {
		l.Error(err)
		return
	}
	if fw, err = os.OpenFile(fullName, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		l.Error(err)
		return
	}
	defer fw.Close()
	if _, err = io.Copy(fw, fr); err != nil {
		l.Error(err)
		return
	}

	sum, err = util.FileMD5(fullName)
	if err != nil {
		l.Error(err)
	}
	return
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsNotExist(err)
}

func IsDirExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsNotExist(err)) && (info != nil && info.IsDir())
}

func UnTar(dst, src string) (err error) {
	if !IsDirExists(dst) {
		if err = os.MkdirAll(dst, 0777); err != nil {
			return
		}
	}

	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	tr := tar.NewReader(fr)
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		filename := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if !IsDirExists(filename) {
				err = os.MkdirAll(filename, 0775)
			}
		case tar.TypeReg:
			err = func() error {
				f, err := os.OpenFile(
					filename, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode),
				)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.Copy(f, tr)
				return err
			}()
		default:
			continue // skip other flag
		}
		if err != nil {
			return err
		}
	}
}

func CheckMD5(filename, sum string) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	_sum := fmt.Sprintf("%x", h.Sum(nil))

	if _sum != sum {
		return errors.New("md5 checksum failed")
	}
	return nil
}
