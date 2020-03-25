package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

/*
	README

	Set next environmental variables:
		AWS_ACCESS_KEY_ID
		AWS_SECRET_ACCESS_KEY
		AWS_REGION=eu-north-1
*/

const (
	S3Region = "eu-north-1"
	S3Bucket = "eventum"
)

type S3Storage struct {
	sess *session.Session
}

func StartAWS() (*S3Storage, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(S3Region),
			Credentials: credentials.NewEnvCredentials(),
		},
	})

	return &S3Storage{sess: sess}, err
}

func (s *S3Storage) UploadToAWS(imgBase64 *string, folder string, name string) error {
	// Config settings: this is where you choose the bucket, filename,
	// content-type etc of the file you're uploading.
	_, err := s3.New(s.sess).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(S3Bucket),
		Key:                aws.String(folder + "/" + name),
		Body:               bytes.NewReader(bytes.NewBufferString(*imgBase64).Bytes()),
		ACL:                aws.String("private"),
		ContentDisposition: aws.String("attachment"),
	})

	return err
}

func (s *S3Storage) ListObjects(folder string, limit int64) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(S3Bucket),
		MaxKeys: aws.Int64(limit),
		Prefix:  aws.String(folder + "/"),
	}

	result, err := s3.New(s.sess).ListObjectsV2(input)

	if err != nil {
		return nil, err
	}

	return result, err
}
