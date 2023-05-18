package resource

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"io"
	"mime/multipart"
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/types"
)

var reserve = int64(100 * 1024 * 1024)

func CheckFileMd5SumAndGetData(ctx context.Context, fh *multipart.FileHeader, md5Str string) (data []byte, sum string, err error) {
	uploadConf := types.MustUploadConfigFromContext(ctx)

	limit := uploadConf.FilesizeLimitBytes
	diskReserve := uploadConf.DiskReserveBytes

	if diskReserve != 0 {
		info, _err := disk.Usage(os.TempDir())
		if _err != nil {
			err = status.UploadFileFailed.StatusErr().WithDesc(_err.Error())
			return
		}
		if info.Free < uint64(diskReserve) {
			err = status.UploadFileDiskLimit
			return
		}
	}

	f, _err := fh.Open()
	if _err != nil {
		err = status.UploadFileFailed.StatusErr().WithDesc(_err.Error())
		return
	}
	defer f.Close()

	data, err = io.ReadAll(f)
	if err != nil {
		return
	}

	if limit > 0 {
		if int64(len(data)) > limit {
			err = status.UploadFileSizeLimit
			return
		}
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
	fs := types.MustFileSystemOpFromContext(ctx)

	path = fmt.Sprintf("%s/%d", os.Getenv(consts.EnvResourceGroup), id)
	err = fs.Upload(path, data)
	if err != nil {
		err = status.UploadFileFailed.StatusErr().WithDesc(err.Error())
		return
	}
	return
}
