package aws

import (
	"bytes"
	"failless/internal/pkg/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"log"
)

const (
	S3Region = "eu-north-1"
	S3Bucket = "https://eventum.s3.eu-north-1.amazonaws.com"
)

func startAWS() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(S3Region),
			Endpoint: aws.String(S3Bucket),
		},
	})

	if err != nil {
		log.Printf("Unable to start session %v", err)
		return nil
	}

	return sess
}

func UploadToAWS(s *session.Session, pic *forms.EImage, folder string) error {
	// Config settings: this is where you choose the bucket, filename,
	// content-type etc of the file you're uploading.
    numBytes, err := s3.New(s).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3Bucket),
        Key:                  aws.String("/" + folder + "/" + pic.ImgName),
        Body:                 bytes.NewReader(bytes.NewBufferString(pic.ImgBase64).Bytes()),
        ACL: 				  aws.String("private"),
        ContentDisposition:   aws.String("attachment"),
    })

    if err != nil {
    	log.Printf("Unable to upload item %q, %v", pic.ImgName, err)
        return err
    } else {
    	log.Printf("Uploaded %q bytes", numBytes)
	}

    return err
}

func DownloadFromAWS(s *session.Session, pic *forms.EImage, folder string) error {
	downloader := s3manager.NewDownloader(s)

	buf := aws.NewWriteAtBuffer([]byte{})
    numBytes, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String("/" + folder + "/" + pic.ImgName),
	})

    pic.Img, _, _ = image.Decode(bytes.NewReader(buf.Bytes()))

    if err != nil {
        log.Printf("Unable to download item %q, %v", pic.ImgName, err)
        return err
    } else {
    	log.Printf("Donwloaded %q bytes", numBytes)
	}

    return err
}
