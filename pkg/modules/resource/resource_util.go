package resource

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/types"
)

var (
	Upload = uploadLocal
	Load   = loadFromLocal
	// Upload = uploadS3
	// Load   = loadFromS3
)

var reserve = int64(100 * 1024 * 1024)

func uploadLocal(ctx context.Context, name string, f io.ReadSeeker) (filename string, sum string, err error) {
	c := types.MustUploadConfigFromContext(ctx)

	if err = checkDisk(ctx, c, f); err != nil {
		return
	}

	filename = filepath.Join(c.Root, name)

	var ff *os.File
	if ff, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		return
	}
	defer ff.Close()
	if _, err = io.Copy(ff, f); err != nil {
		return
	}

	sum, err = util.FileMD5(filename)
	return
}

func loadFromLocal(ctx context.Context) ([]byte, error) {
	res := types.MustResourceFromContext(ctx)
	code, err := os.ReadFile(res.Path)
	if err != nil {
		return nil, status.LoadLocalWasmFailed.StatusErr().WithDesc(err.Error())
	}
	return code, nil
}

func uploadS3(ctx context.Context, name string, f io.ReadSeeker) (uri string, sum string, err error) {
	// TODO
	return
}

func loadFromS3(ctx context.Context) (raw []byte, err error) {
	// TODO
	return
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsNotExist(err)) && (info != nil && info.IsDir())
}

func checkDisk(ctx context.Context, c *types.UploadConfig, f io.ReadSeeker) (err error) {
	var (
		size = int64(0)
		stat *disk.UsageStat
	)

	if !isDirExists(c.Root) {
		if err = os.MkdirAll(c.Root, 0777); err != nil {
			return
		}
	}

	if size, err = f.Seek(0, io.SeekEnd); err != nil {
		return
	}
	if size > c.FileSizeLimit {
		err = errors.Errorf("filesize over limit")
		return
	}

	stat, err = disk.Usage(c.Root)
	if err != nil {
		return
	}

	if stat.Free < uint64(size+reserve) {
		err = errors.Wrap(err, "disk limited")
		return
	}
	_, err = f.Seek(0, io.SeekStart)
	return err
}
