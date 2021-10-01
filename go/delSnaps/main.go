package main

import (
    "context"
    "fmt"
    "log"

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

    // get snapshot by description value
    snapsRes, err := svc.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
        Filters: []types.Filter{
            {
                Name: aws.String("description"),
                Values: []string{"Created by Lambda backup function ebs-snapshots"},
            },
        },
    })
    if err != nil {
        log.Fatalf("unable to DescribeSnapshots with filter, %v", err)
    }

    total := len(snapsRes.Snapshots)
    fmt.Println("got", total, "snapshots");

    // delete all the old snapshots
    for index, snapShot := range snapsRes.Snapshots {
        fmt.Println("deleteing", *snapShot.SnapshotId, index + 1, "of", total)
        delSnap, err := svc.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
            SnapshotId: aws.String(*snapShot.SnapshotId),
            DryRun: aws.Bool(false),
        })
        if err != nil {
            log.Fatalf("unable to delete snapshot: %s, Error: %v", *snapShot.SnapshotId, err)
        }
        fmt.Println("Return Metadata:", *delSnap)
    }
}
