package resource

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"io"
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v3/disk"

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

func CheckFileMd5Sum(f io.ReadSeekCloser, md5Str string) (data []byte, sum string, err error) {
	data, err = io.ReadAll(f)
	if err != nil {
		return
	}
	hash := md5.New()
	_, err = hash.Write(data)
	if err != nil {
		return
	}

	sum = fmt.Sprintf("%x", hash.Sum(nil))
	if md5Str != "" && sum != md5Str {
		err = status.UploadFileMd5Unmatched
		return
	}
	return
}

func UploadFile(ctx context.Context, data []byte, id types.SFID) (path string, err error) {
	var (
		fs         = types.MustFileSystemOpFromContext(ctx)
		uploadConf = types.MustUploadConfigFromContext(ctx)
	)

	limit := uploadConf.FilesizeLimitBytes
	diskReserve := uploadConf.DiskReserveBytes
	root := uploadConf.Root

	if limit > 0 {
		if int64(len(data)) > limit {
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

	path = fmt.Sprintf("%s/%d", os.Getenv(consts.EnvResourceGroup), id)
	err = fs.Upload(path, data)
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
