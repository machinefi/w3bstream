package resource

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/local"
	"github.com/machinefi/w3bstream/pkg/errors/status"
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
		fs          = types.MustFileSystemOpFromContext(ctx)
		size        = int64(0)
		limit       = int64(0)
		diskReserve = int64(0)
		root        = ""
	)
	if v, ok := fs.(*local.LocalFileSystem); ok {
		limit = v.FilesizeLimitBytes
		diskReserve = v.DiskReserveBytes
		root = v.Root
	}

	if limit > 0 {
		if err, _ = checkFilesize(f, limit); err != nil {
			if err != nil {
				err = status.UploadFileFailed.StatusErr().WithDesc(err.Error())
				return
			}
		}
		if size > limit {
			err = status.UploadFileSizeLimit
			return
		}
	}

	if root != "" && diskReserve != 0 {
		info, _err := disk.Usage(root)
		if _err != nil {
			err = status.UploadFileFailed.StatusErr().WithDesc(_err.Error())
			return
		}
		if info.Free < uint64(diskReserve) {
			err = status.UploadFileDiskLimit
			return
		}
	}

	sum := ""
	data, sum, err = checkFileMd5Sum(f)
	if err != nil {
		err = status.MD5ChecksumFailed
		return
	}

	if md5 != "" && sum != md5 {
		err = status.UploadFileMd5Unmatched
		return
	}

	path = md5
	err = fs.Upload(md5, data)
	if err != nil {
		err = status.UploadFileFailed.StatusErr().WithDesc(err.Error())
		return
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
