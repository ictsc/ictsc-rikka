package s3repo

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

var (
	bucketname = "ictsc"
)

type S3Repository struct {
	minioclient *minio.Client
}

func NewS3Repository(minioclient *minio.Client) *S3Repository {
	return &S3Repository{
		minioclient: minioclient,
	}
}

func (r *S3Repository) Create(id string, reader io.Reader) error {
	err := r.minioclient.MakeBucket(context.Background(), bucketname, minio.MakeBucketOptions{ObjectLocking: true})
	if err != nil {
		exists, errBucketExists := r.minioclient.BucketExists(context.Background(), bucketname)
		if errBucketExists == nil && exists {
			_, errPut := r.minioclient.PutObject(context.Background(), bucketname, id, reader, -1, minio.PutObjectOptions{})
			if errPut != nil {
				return err
			}
			return nil
		}
		return err
	}
	_, errPut := r.minioclient.PutObject(context.Background(), bucketname, id, reader, -1, minio.PutObjectOptions{})
	if errPut != nil {
		return err
	}
	return nil

}
func (r *S3Repository) Delete(id string) error {
	err := r.minioclient.RemoveObject(context.Background(), bucketname, id, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil

}
func (r *S3Repository) Get(id string) (io.Reader, error) {
	obj, err := r.minioclient.GetObject(context.Background(), bucketname, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil

}
func (r *S3Repository) GetAll() ([]*minio.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objCh := r.minioclient.ListObjects(ctx, bucketname, minio.ListObjectsOptions{Recursive: true})
	var result []*minio.ObjectInfo
	for obj := range objCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		result = append(result, &obj)
	}
	return result, nil

}
