package aws

import (
    "failless/internal/pkg/forms"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func UploadToAWS(pic *forms.EImage) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1")},
	)
}
