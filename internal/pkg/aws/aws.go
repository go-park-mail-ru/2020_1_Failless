package aws

import (
	"bytes"
	"failless/internal/pkg/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
)

const (
	S3Region = "eu-north-1"
	S3Bucket = "eventum"
	// AWS_PROFILE=test-account
)

func StartAWS() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(S3Region),
			Credentials: credentials.NewSharedCredentials("", "test-account"),
		},
	})

	return sess, err
}

func UploadToAWS(sess *session.Session, pic *forms.EImage, folder string) error {
	// Config settings: this is where you choose the bucket, filename,
	// content-type etc of the file you're uploading.
    _, err := s3.New(sess).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3Bucket),
        Key:                  aws.String(folder + "/" + pic.ImgName),
        Body:                 bytes.NewReader(bytes.NewBufferString(pic.ImgBase64).Bytes()),
        ACL: 				  aws.String("private"),
        ContentDisposition:   aws.String("attachment"),
    })

    return err
}

func DownloadFromAWS(sess *session.Session, pic *forms.EImage, folder string) error {
	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})
    _, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(folder + "/" + pic.ImgName),
	})

    pic.Img, _, err = image.Decode(bytes.NewReader(buf.Bytes()))

    return err
}

func ListObjects(sess *session.Session, folder string, limit int64) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(S3Bucket),
		MaxKeys: aws.Int64(limit),
		Prefix:  aws.String(folder + "/"),
	}

	result, err := s3.New(sess).ListObjectsV2(input)

	if err != nil {
		return nil, err
	}

	return result, err
}
