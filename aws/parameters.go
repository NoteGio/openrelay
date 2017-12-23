package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParameter(parameter string) (string) {
	sess := session.Must(session.NewSession())
	client := ssm.New(sess)
	paramValue, _ := client.GetParameter(&ssm.GetParameterInput{Name: &parameter, WithDecryption: aws.Bool(true)})
	return paramValue.GoString()
}
