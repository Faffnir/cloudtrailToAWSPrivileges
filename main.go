package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type awsCloudJSONObject struct {
	Records []awsCloudJSONElement `json:"Records"`
}

type awsCloudJSONElement struct {
	EventName   string `json:"eventName"`
	EventSource string `json:"eventSource"`
}

type awsPermissions struct {
	Version   string                 `json:"Version"`
	Statement []awsPermissionElement `json:"Statement"`
}

type awsPermissionElement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

func main() {
	awsCloudTrailFile := flag.String("file", "", "The AWS Cloudtrail json file")

	flag.Parse()

	if *awsCloudTrailFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*awsCloudTrailFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jsonByteArray, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var data awsCloudJSONObject

	err = json.Unmarshal(jsonByteArray, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Deduplicate permissions
	rules := map[string]bool{}
	for _, element := range data.Records {
		eventSource := strings.Split(element.EventSource, ".")[0]
		rules[eventSource+":"+element.EventName] = true
	}

	// Create a slice out of all permissions
	keys := make([]string, 0, len(rules))
	for k := range rules {
		keys = append(keys, k)
	}

	awsPermissionBlob := awsPermissions{
		Version: "2012-10-17",
		Statement: []awsPermissionElement{
			awsPermissionElement{
				Effect:   "Allow",
				Action:   keys,
				Resource: "*",
			},
		},
	}

	permissionString, err := json.Marshal(awsPermissionBlob)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", string(permissionString))
}
