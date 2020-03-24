package aws

import (
	"failless/internal/pkg/forms"
	"github.com/stretchr/testify/assert"
	"os"

	//"failless/internal/pkg/forms"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	//"os"
	"testing"
)

const (
	testImageName = "default.png"
	testFolder = "test"
)

func TestStartAWS(t *testing.T) {
	sess, err := StartAWS()
	if err != nil {
		log.Printf("Unable to start session %v", err)
		t.Fail()
	} else {
		creds, err := sess.Config.Credentials.Get()
		if err != nil {
			log.Printf("Unable to get creds %v", err)
			t.Fail()
		} else {
			log.Printf("Creds: %v", creds)
		}
	}
}

func TestListObjects(t *testing.T) {
	sess, _ := StartAWS()
	_, err := ListObjects(sess, "event", 10)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println("Answer Error: ", aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		t.Fail()
	}
}

func TestUploadToAWS(t *testing.T) {
	path := "../" + forms.Media + testImageName
	sess, _ := StartAWS()

	_, err := os.Stat(path)
	if err != nil {
		log.Printf("File '%q' doesn't exist\n%v", testImageName, err)
		if pwd, err := os.Getwd(); err != nil {
			log.Printf("os.Getwd() failed with %v", err)
			t.Fail()
		} else {
			log.Printf("pwd: %q", pwd)
		}
		t.Fail()
	}

	tempImg := correctImage(path, testImageName)

	err = UploadToAWS(sess, &tempImg, testFolder)
	if err != nil {
		log.Printf("Unable to upload item %q, %v", tempImg.ImgName, err)
		t.Fail()
	} else {
		result, _ := ListObjects(sess, testFolder, 2)
		if len(result.Contents) < 2 {
			log.Println("Impossible.\nPerhaps the archives are incomplete.")
			t.Fail()
		}
		assert.Equal(t, testFolder + "/" + testImageName, *result.Contents[1].Key)
	}
}

func correctImage(path, name string) (tmp forms.EImage) {
	_ = tmp.GetImage(path)
	_ = tmp.Encode()
	tmp.ImgName = name
	return
}
