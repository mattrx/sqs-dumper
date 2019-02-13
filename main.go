package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

var (
	profile           string
	queue             string
	output            string
	loopCount         int32
	visibilityTimeout int32
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "./messages", "output directory")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "aws profile")
	rootCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "queue url")
	rootCmd.PersistentFlags().Int32VarP(&loopCount, "loop-count", "", 1000, "number of loops for receive")
	rootCmd.PersistentFlags().Int32VarP(&visibilityTimeout, "visibility-timeout", "", 60, "visibility timeout in seconds")
}

var rootCmd = &cobra.Command{
	Use: "sqs-dumper",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(output); os.IsNotExist(err) {
			if err := os.Mkdir(output, 0777); err != nil {
				return err
			}
		}

		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           profile,
		})

		if err != nil {
			return err
		}

		svc := sqs.New(sess)

		for i := 1; i <= int(loopCount); i++ {
			receiveResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queue),
				MaxNumberOfMessages: aws.Int64(10),
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				WaitTimeSeconds:   aws.Int64(10),
				VisibilityTimeout: aws.Int64(int64(visibilityTimeout)),
			})

			if err != nil {
				return err
			}

			for _, message := range receiveResult.Messages {
				path := strings.TrimRight(output, "/") + "/" + *message.MessageId + ".json"
				if _, err := os.Stat(path); os.IsNotExist(err) {
					if err := ioutil.WriteFile(path, []byte(*message.Body), 0644); err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}
