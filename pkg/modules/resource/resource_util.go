package resource

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/types"
)

var reserve = int64(100 * 1024 * 1024)

func checkFilesize(f io.ReadSeekCloser, lmt int64) (err error, size int64) {
	size, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	return nil, size
}

func checkFileMd5Sum(f io.Reader) (data []byte, sum string, err error) {
	data, err = io.ReadAll(f)
	if err != nil {
		return
	}
	hash := md5.New()
	_, err = hash.Write(data)
	if err != nil {
		return
	}

	return data, fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func UploadFile(ctx context.Context, f io.ReadSeekCloser, md5 string) (path string, data []byte, err error) {
	var (
		c    = types.MustUploadConfigFromContext(ctx)
		size = int64(0)
	)

	if err, size = checkFilesize(f, c.FileSizeLimit); err != nil {
		if err != nil {
			return
		}
	}

	sum := ""
	data, sum, err = checkFileMd5Sum(f)
	if err != nil {
		return
	}

	if sum != md5 {
		err = errors.Errorf("file md5 not match")
		return
	}

	path = md5
	if fs, ok := types.FileSystemFromContext(ctx); ok && fs != nil {
		// store to s3
		err = fs.FileCli.Upload(path, data)
	} else {
		// store local
		var (
			lf    *os.File
			stat  *disk.UsageStat
			wsize int
		)
		stat, err = disk.Usage(c.Root)
		if err != nil {
			return
		}
		if stat == nil || stat.Free < uint64(c.FileSizeLimit+reserve) {
			err = errors.New("disk limited")
			return
		}
		if !IsDirExists(c.Root) {
			if err = os.MkdirAll(c.Root, 0777); err != nil {
				return
			}
		}
		if lf, err = os.Create(filepath.Join(c.Root, md5)); err != nil {
			return
		}
		defer lf.Close()
		wsize, err = lf.Write(data)
		if err != nil {
			return
		}
		if int64(wsize) != size {
			err = errors.New("write file failed")
			return
		}
		err = os.Chmod(filepath.Join(c.Root, md5), 0400)
	}
	return
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
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
