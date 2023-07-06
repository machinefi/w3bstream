package confs3

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
)

type S3Endpoint interface {
	Endpoint() string
	AccessKeyID() string
	SecretAccessKey() string
	BucketName() string
	Secure() bool
}

type PresignedFn func(db *ObjectDB, key string, exp time.Duration) url.Values

// ObjectDB Deprecated
type ObjectDB struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey types.Password
	BucketName      string
	Secure          bool
	UrlExpire       types.Duration
	Presigned       PresignedFn `env:"-"`
}

func (db *ObjectDB) SetDefault() {
	if db.UrlExpire == 0 {
		db.UrlExpire = types.Duration(10 * time.Minute)
	}
}

func (db *ObjectDB) LivenessCheck() map[string]string {
	key := db.BucketName + "." + db.Endpoint
	m := map[string]string{
		key: "ok",
	}

	c, err := db.Client()

	if err != nil {
		m[key] = err.Error()
	} else {
		if _, err := c.GetBucketLocation(context.Background(), db.BucketName); err != nil {
			m[key] = err.Error()
		}
	}

	return m
}

func (db *ObjectDB) Client() (*minio.Client, error) {
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(db.AccessKeyID, db.SecretAccessKey.String(), ""),
		Secure: db.Secure,
		Region: db.Region,
	}

	client, err := minio.New(db.Endpoint, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (db *ObjectDB) PublicURL(meta *filesystem.ObjectMeta) *url.URL {
	u := &url.URL{}
	u.Scheme = "http"
	if db.Secure {
		u.Scheme += "s"
	}

	u.Host = db.Endpoint
	u.Path = db.BucketName + "/" + meta.Key()
	return u
}

func (db *ObjectDB) ProtectURL(ctx context.Context, meta *filesystem.ObjectMeta, exp time.Duration) (*url.URL, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	if db.Presigned != nil {
		values = db.Presigned(db, meta.Key(), exp)
	}

	u, err := c.PresignedGetObject(ctx, db.BucketName, meta.Key(), exp, values)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *ObjectDB) PutObject(ctx context.Context, r io.Reader, meta *filesystem.ObjectMeta) error {
	if ctx == nil {
		ctx = context.Background()
	}

	c, err := db.Client()
	if err != nil {
		return err
	}

	if meta.Size == 0 {
		if canLen, ok := r.(interface{ Len() int }); ok {
			meta.Size = int64(canLen.Len())
		}
	}

	_, err = c.PutObject(
		ctx, db.BucketName, meta.Key(),
		r, meta.Size, minio.PutObjectOptions{ContentType: meta.ContentType},
	)

	return err
}

func (db *ObjectDB) ReadObject(ctx context.Context, w io.Writer, meta *filesystem.ObjectMeta) error {
	c, err := db.Client()
	if err != nil {
		return err
	}

	object, err := c.GetObject(ctx, db.BucketName, meta.Key(), DefaultGetObjectOptions)
	if err != nil {
		return err
	}
	defer object.Close()

	_, err = io.Copy(w, object)
	if err != nil {
		return err
	}

	return err
}

func (db *ObjectDB) PresignedPutObject(ctx context.Context, meta *filesystem.ObjectMeta, exp time.Duration) (string, error) {
	c, err := db.Client()
	if err != nil {
		return "", err
	}
	address, err := c.PresignedPutObject(
		ctx,
		db.BucketName,
		meta.Key(),
		exp,
	)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (db *ObjectDB) DeleteObject(ctx context.Context, meta *filesystem.ObjectMeta) error {
	c, err := db.Client()
	if err != nil {
		return err
	}

	return c.RemoveObject(ctx, db.BucketName, meta.Key(), DefaultRemoveObjectOptions)
}

func (db *ObjectDB) StatsObject(ctx context.Context, meta *filesystem.ObjectMeta) (*filesystem.ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}

	info, err := c.StatObject(ctx, db.BucketName, meta.Key(), minio.StatObjectOptions{})
	if err != nil {
		awsErr, ok := err.(awserr.RequestFailure)
		if ok && awsErr.StatusCode() == 404 {
			return nil, filesystem.ErrNotExistObjectKey
		}
		return nil, err
	}

	om, err := filesystem.ParseObjectMetaFromKey(info.Key)
	if err != nil {
		return nil, err
	}

	om.ContentType = info.ContentType
	om.ETag = info.ETag
	om.Size = info.Size

	return om, err
}

func (db *ObjectDB) ListObjectByGroup(ctx context.Context, grp string) ([]*filesystem.ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}

	metas := make([]*filesystem.ObjectMeta, 0)

	objectsCh := c.ListObjects(ctx, db.BucketName, minio.ListObjectsOptions{
		Prefix:    grp,
		Recursive: true,
	})

	for obj := range objectsCh {
		om, err := filesystem.ParseObjectMetaFromKey(obj.Key)
		if err != nil {
			continue
		}

		om.ContentType = obj.ContentType
		om.ETag = obj.ETag
		om.Size = obj.Size

		metas = append(metas, om)
	}

	return metas, nil
}

func (db *ObjectDB) Upload(key string, content []byte) error {
	meta, err := filesystem.ParseObjectMetaFromKey(key)
	if err != nil {
		return err
	}
	return db.PutObject(context.Background(), bytes.NewBuffer(content), meta)
}

func (db *ObjectDB) Read(key string) ([]byte, error) {
	meta, err := filesystem.ParseObjectMetaFromKey(key)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = db.ReadObject(context.Background(), buf, meta)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (db *ObjectDB) Delete(key string) error {
	meta, err := filesystem.ParseObjectMetaFromKey(key)
	if err != nil {
		return err
	}
	return db.DeleteObject(context.Background(), meta)
}

func (db *ObjectDB) DownloadUrl(key string) (string, error) {
	meta, err := filesystem.ParseObjectMetaFromKey(key)
	if err != nil {
		return "", err
	}
	u, err := db.ProtectURL(context.Background(), meta, db.UrlExpire.Duration())
	if err != nil {
		return "", err
	}
	return u.String(), err
}

// StatObject Deprecated
func (db *ObjectDB) StatObject(key string) (*filesystem.ObjectMeta, error) {
	meta, err := filesystem.ParseObjectMetaFromKey(key)
	if err != nil {
		return nil, err
	}
	return db.StatsObject(context.Background(), meta)
}

var (
	DefaultGetObjectOptions    = minio.GetObjectOptions{}
	DefaultRemoveObjectOptions = minio.RemoveObjectOptions{}
)
