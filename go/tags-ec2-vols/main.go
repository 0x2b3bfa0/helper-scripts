package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"

	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the ec2 client
	svc := ec2.NewFromConfig(cfg)

	// get volumes
	volumeRes, err := svc.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(1000),
	})
	if err != nil {
		log.Fatal(err)
	}

	for /*volIndex*/ _, volume := range volumeRes.Volumes {
		// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2@v1.18.0/types#Volume
		//volumeId := *volume.VolumeId
		volTags := volume.Tags
		volAttachment := volume.Attachments

		for /*volAi*/ _, attachment := range volAttachment {
			instanceId := *attachment.InstanceId
			instancesRes, _ := svc.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
				InstanceIds: []string{instanceId},
			})
			if len(instancesRes.Reservations) > 1 {
				log.Fatalf("instance resverations bigger than expected > 1")
			}
			//fmt.Println(len(instancesRes.Reservations))
			//fmt.Println(instancesRes.Reservations[0])
			///*
			if len(instancesRes.Reservations[0].Instances) > 1 {
				log.Fatalf("only expected to get 1 instance");
			}
			///*
			instance := instancesRes.Reservations[0].Instances[0]
			instanceTags := instance.Tags
			fmt.Println("Instance: " + *instance.InstanceId + " has tags:")
			for _, tag := range instanceTags{
				fmt.Println(*tag.Key + " = " + *tag.Value)
			}
			fmt.Println()
			missingTags(instanceTags, volTags)
			fmt.Println("-----------------")
			//fmt.Println(len(instanceTags))
			//*/
		}
	}
}

// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2@v1.18.0#Client.CreateTags
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2@v1.18.0#CreateTagsInput
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2@v1.18.0#CreateTagsOutput


func missingTags(iTags []types.Tag, vTags []types.Tag) { //[]Types.Tag {
	if len(iTags) == 0 {
		fmt.Println("WARNING no Tags on Instance")
		return
	}
	if len(vTags) == 0 {
		fmt.Println("WARNING no Tags on Volume")
		return
	}
	for _, iTag := range iTags {
		for vi, vTag := range vTags {
			if *iTag.Key == *vTag.Key {
				// has key
				if *iTag.Value == *vTag.Value {
					// values match
					fmt.Println("Volume has key:", *iTag.Key, "with same value")
					break
				} else {
					// values dont match
					fmt.Println("Volume has key:", *iTag.Key, "with different value")
					// handle case
					break
				}
			}
			if vi == len(vTags) - 1 {
				// on last key and not found
				fmt.Println("Volume doesnt have key: " + *iTag.Key)
			}
		}
	}
}

/*
	// get snapshot by description value
	instancesRes, err := svc.DescribeInstances(ctx, &ec2.DescribeSnapshotsInput{
		MaxResults: aws.Int32(1000)
	})
	if err != nil {
		log.Fatalf("Error svr.DescribeInstances, %v", err)
	}

	for index, reservation := range instancesRes.Reservations {
		for i, instance := range reservation.Instance {
			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2@v1.18.0/types#Instance
			ec2tags := instance.Tags
		}
	}
}
*/
