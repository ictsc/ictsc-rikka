package s3repo

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type S3Repository struct {
	minioclient *minio.Client
	bucketName  string
}

func NewS3Repository(minioclient *minio.Client, bucketName string) *S3Repository {
	return &S3Repository{
		minioclient: minioclient,
		bucketName:  bucketName,
	}
}

func (r *S3Repository) Create(id string, reader io.Reader) error {
	_, err := r.minioclient.PutObject(context.Background(), r.bucketName, id, reader, -1, minio.PutObjectOptions{})
	return err
}
func (r *S3Repository) Delete(id string) error {
	err := r.minioclient.RemoveObject(context.Background(), r.bucketName, id, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (r *S3Repository) Get(id string) (io.Reader, error) {
	obj, err := r.minioclient.GetObject(context.Background(), r.bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil

}
func (r *S3Repository) GetAll() ([]*minio.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objCh := r.minioclient.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{Recursive: true})
	var result []*minio.ObjectInfo
	for obj := range objCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		result = append(result, &obj)
	}
	return result, nil

}
