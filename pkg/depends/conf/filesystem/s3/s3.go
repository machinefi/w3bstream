package confs3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type S3Endpoint interface {
	Endpoint() string
	AccessKeyID() string
	SecretAccessKey() string
	BucketName() string
	Secure() bool
}

type PresignedFn func(db *ObjectDB, key string, exp time.Duration) url.Values

type ObjectDB struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey types.Password
	BucketName      string
	Secure          bool
	Presigned       PresignedFn `env:"-"`
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

func (db *ObjectDB) PublicURL(meta *ObjectMeta) *url.URL {
	u := &url.URL{}
	u.Scheme = "http"
	if db.Secure {
		u.Scheme += "s"
	}

	u.Host = db.Endpoint
	u.Path = db.BucketName + "/" + meta.Key()
	return u
}

func (db *ObjectDB) ProtectURL(ctx context.Context, meta *ObjectMeta, exp time.Duration) (*url.URL, error) {
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

func (db *ObjectDB) PutObject(ctx context.Context, r io.Reader, meta *ObjectMeta) error {
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

func (db *ObjectDB) ReadObject(ctx context.Context, w io.Writer, grp string, oid uint64) error {
	c, err := db.Client()
	if err != nil {
		return err
	}
	key := (&ObjectMeta{Group: grp, ObjectID: oid}).Key()

	object, err := c.GetObject(ctx, db.BucketName, key, DefaultGetObjectOptions)
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

func (db *ObjectDB) PresignedPutObject(ctx context.Context, grp string, oid uint64, exp time.Duration) (string, error) {
	c, err := db.Client()
	if err != nil {
		return "", err
	}
	address, err := c.PresignedPutObject(
		ctx,
		db.BucketName,
		(&ObjectMeta{Group: grp, ObjectID: oid}).Key(),
		exp,
	)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (db *ObjectDB) DeleteObject(ctx context.Context, grp string, oid uint64) error {
	c, err := db.Client()
	if err != nil {
		return err
	}

	key := (&ObjectMeta{Group: grp, ObjectID: oid}).Key()

	return c.RemoveObject(ctx, db.BucketName, key, DefaultRemoveObjectOptions)
}

func (db *ObjectDB) StatsObject(ctx context.Context, grp string, oid uint64) (*ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}
	key := (&ObjectMeta{Group: grp, ObjectID: oid}).Key()

	object, err := c.GetObject(ctx, db.BucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	info, err := object.Stat()
	if err != nil {
		return nil, err
	}

	om, err := ParseObjectMetaFromKey(info.Key)
	if err != nil {
		return nil, err
	}

	om.ContentType = info.ContentType
	om.ETag = info.ETag
	om.Size = info.Size

	return om, err
}

func (db *ObjectDB) ListObjectByGroup(ctx context.Context, grp string) ([]*ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}

	metas := make([]*ObjectMeta, 0)

	objectsCh := c.ListObjects(ctx, db.BucketName, minio.ListObjectsOptions{
		Prefix:    grp,
		Recursive: true,
	})

	for obj := range objectsCh {
		om, err := ParseObjectMetaFromKey(obj.Key)
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
	meta, err := ParseObjectMetaFromKey(key)
	if err != nil {
		return err
	}
	return db.PutObject(context.Background(), bytes.NewBuffer(content), meta)
}

func (db *ObjectDB) Read(key string) ([]byte, error) {
	meta, err := ParseObjectMetaFromKey(key)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = db.ReadObject(context.Background(), buf, meta.Group, meta.ObjectID)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (db *ObjectDB) Delete(key string) error {
	meta, err := ParseObjectMetaFromKey(key)
	if err != nil {
		return err
	}
	return db.DeleteObject(context.Background(), meta.Group, meta.ObjectID)
}

var ErrInvalidObjectKey = errors.New("invalid object key")

func ParseObjectMetaFromKey(key string) (*ObjectMeta, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidObjectKey
	}
	grp := parts[0]

	oid, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidObjectKey
	}

	om := &ObjectMeta{ObjectID: oid, Group: grp}

	return om, nil
}

var (
	DefaultGetObjectOptions    = minio.GetObjectOptions{}
	DefaultRemoveObjectOptions = minio.RemoveObjectOptions{}
)
