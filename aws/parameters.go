package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
)

func GetParameter(parameter string) string {
	sess := session.Must(session.NewSession())
	client := ssm.New(sess)
	paramValue, err := client.GetParameter(&ssm.GetParameterInput{Name: &parameter, WithDecryption: aws.Bool(true)})
	if err != nil {
		log.Printf("Error getting Parameter %v: %v", parameter, err.Error())
	}
	return *paramValue.Parameter.Value
}
