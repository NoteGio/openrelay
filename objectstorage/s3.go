package objectstorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3mod "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"bytes"
)

type S3Storage struct {
	bucket string
	key string
}

func (s3 *S3Storage) Reader() (io.Reader, error) {
	sess, _ := session.NewSession(&aws.Config{})
	downloader := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf,
		&s3mod.GetObjectInput{
			Bucket: aws.String(s3.bucket),
			Key:    aws.String(s3.key),
		})
	return bytes.NewBuffer(buf.Bytes()), err
}

func (s3 *S3Storage) Writer() (io.WriteCloser, error) {
	return &S3Writer{s3.bucket, s3.key, &bytes.Buffer{}}, nil
}

type S3Writer struct {
	bucket string
	key string
	buf *bytes.Buffer
}

func (s3w *S3Writer) Write(data []byte) (int, error) {
	return s3w.buf.Write(data)
}

func (s3w *S3Writer) Close() error {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)
	// Upload the file to S3.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3w.bucket),
		Key:    aws.String(s3w.key),
		Body:   s3w.buf,
	})
	if err != nil { return err }
	return nil
}
