package aws

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

const (
	testImageName = "default.png"
	testFolder    = "test"
)

func TestStartAWS(t *testing.T) {
	sess, err := StartAWS()
	if err != nil {
		log.Printf("Unable to start session %v", err)
		t.Fail()
	} else {
		credentials, err := sess.sess.Config.Credentials.Get()
		if err != nil {
			log.Printf("Unable to get credentials %v", err)
			t.Fail()
		} else {
			log.Printf("Credentials: %v", credentials)
		}
	}
}

func TestListObjects(t *testing.T) {
	sess, _ := StartAWS()
	_, err := sess.ListObjects("event", 10)
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
	path := "../../../media/images/" + testImageName
	service, _ := StartAWS()

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
		return
	}

	buf, err := correctImage(path, testImageName)
	if err != nil {
		t.Fail()
		return
	}

	err = service.UploadToAWS(bytes.NewReader(buf.Bytes()), testFolder, testImageName)
	if err != nil {
		log.Printf("Unable to upload item %q, %v", testImageName, err)
		t.Fail()
	} else {
		result, _ := service.ListObjects(testFolder, 2)
		if len(result.Contents) < 2 {
			log.Println("Impossible.\nPerhaps the archives are incomplete.")
			t.Fail()
		}
		assert.Equal(t, testFolder+"/"+testImageName, *result.Contents[1].Key)
	}
}

func correctImage(path, name string) (*bytes.Buffer, error) {
	img, err := imaging.Open(path + name)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
