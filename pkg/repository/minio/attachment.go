package miniorepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/minio/minio-go/v7"
)

var (
	bucketname = "ictsc"
)

type AttachmentRepository struct {
	minioclient *minio.Client
}

func NewAttachmentRepository(minioclient *minio.Client) *AttachmentRepository {
	return &AttachmentRepository{
		minioclient: minioclient,
	}
}

func (r *AttachmentRepository) Upload(attachment *entity.Attachment) (*entity.Attachment, error) {
	err := r.minioclient.MakeBucket(context.Background(), bucketname, minio.MakeBucketOptions{ObjectLocking: true})
	if err != nil {
		exists, errBucketExists := r.minioclient.BucketExists(context.Background(), bucketname)
		if errBucketExists == nil && exists {
			p := make([]byte, 128*1000*1000)
			size, err := attachment.Reader.Read(p)
			if err != nil {
				return nil, err
			}
			//_, errPut := r.minioclient.PutObject(context.Background(), bucketname, attachment.Name)
			_, errPut := r.minioclient.PutObject(context.Background(), bucketname, attachment.Name, attachment.Reader, int64(size), minio.PutObjectOptions{})
			if errPut != nil {
				return nil, err
			}
			return attachment, nil
		}
		return nil, err
	}
	_, errPut := r.minioclient.FPutObject(context.Background(), bucketname, attachment.Name, "", minio.PutObjectOptions{})
	if errPut != nil {
		return nil, err
	}
	return attachment, nil

}
func (r *AttachmentRepository) Delete(attachment *entity.Attachment) error {
	err := r.minioclient.RemoveObject(context.Background(), bucketname, attachment.Name, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil

}
func (r *AttachmentRepository) Get(id uuid.UUID) (*minio.Object, error) {
	obj, err := r.minioclient.GetObject(context.Background(), bucketname, id.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil

}
func (r *AttachmentRepository) GetAll() ([]*minio.ObjectInfo, error) {
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
