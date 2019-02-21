package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	region := flag.String("region", "us-west-1", "AWS Region i.e. us-west-1")
	fil := flag.String("filter", "instance-id", "Filter by name i.e. tag-key, list can be found here https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html")
	value := flag.String("value", "*", "Value of filter i.e. Owner")

	flag.Parse()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*region)},
	)

	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(*fil),
				Values: []*string{
					aws.String(*value),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)
	defer w.Flush()

	//format which includes the instance id, tag value, instance type and launch time
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t", "id", "instance type", "launch time", "tag-key", "tag-value")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t", "---", "---", "---", "---", "---")

	//display any instances without an Owner tag as type 'unknown' with the instance id, type and launch time.
	for i := range result.Reservations {
		for _, inst := range result.Reservations[i].Instances {
			id := *inst.InstanceId
			itype := *inst.InstanceType
			lt := *inst.LaunchTime

			for _, tag := range inst.Tags {
				tagKey := *tag.Key
				tagValue := *tag.Value
				if *tag.Value == "" {
					tagValue = "unknown"
				}
				fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t", id, itype, lt, tagKey, tagValue)
			}

		}
	}
}
