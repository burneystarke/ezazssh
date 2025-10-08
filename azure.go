package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

func GetAzureCred() (*azidentity.AzureCLICredential, error) {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credentials: %w", err)
	}
	return cred, nil
}

func GetSubscriptions() ([]SubscriptionInfo, error) {
	cred, err := GetAzureCred()
	if err != nil {
		return nil, err
	}
	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription client: %w", err)
	}

	pager := client.NewListPager(nil)
	subscriptions := []SubscriptionInfo{}

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get next subscription page: %w", err)
		}
		for _, sub := range page.Value {
			subscriptions = append(subscriptions, SubscriptionInfo{
				ID:   *sub.SubscriptionID,
				Name: *sub.DisplayName,
			})
		}
	}
	return subscriptions, nil
}

// GetMachines retrieves both VMs and ARC-enabled machines using Azure Resource Graph
func GetMachines(subscriptionID string) ([]MachineInfo, error) {
	cred, err := GetAzureCred()
	if err != nil {
		return nil, err
	}

	client, err := armresourcegraph.NewClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource Graph client: %w", err)
	}

	// Azure Resource Graph query
	query := "Resources | where type == 'microsoft.hybridcompute/machines' or type == 'microsoft.compute/virtualmachines' | project id,subscriptionId,resourceGroup, JoinID = toupper(id), ComputerName = name,MachineType=type | join kind=inner( Resources | where type == 'microsoft.hybridcompute/machines/extensions' or type == 'microsoft.compute/virtualmachines/extensions' | where name startswith 'AADSSHLogin' or name startswith 'WindowsOpenSSH' | project  Extension=name, MachineId = toupper(substring(id, 0, indexof(id, '/extensions')))) on $left.JoinID == $right.MachineId | project name=ComputerName,Extension,MachineType, resourceGroup,subscriptionId"
	// Define the query request
	queryRequest := armresourcegraph.QueryRequest{
		Query: to.Ptr(query),
		Subscriptions: []*string{
			to.Ptr(subscriptionID),
		},
	}

	// Execute the query
	results, err := client.Resources(context.Background(), queryRequest, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}

	// Parse the results
	machines := []MachineInfo{}
	rows := results.Data.([]interface{}) // Cast Data to a slice of interfaces
	for _, row := range rows {
		data := row.(map[string]interface{}) // Each row is a map
		machine := MachineInfo{
			Name:           data["name"].(string),
			ResourceGroup:  data["resourceGroup"].(string),
			SubscriptionID: data["subscriptionId"].(string),
		}
		if data["Extension"].(string) == "WindowsOpenSSH" {
			machine.OS = "windows"
		} else {
			machine.OS = "linux"
		}
		// Determine if it's an ARC or Azure VM based on the "type" field
		if data["MachineType"].(string) == "microsoft.hybridcompute/machines" {
			machine.IsArc = true
		} else if data["MachineType"].(string) == "microsoft.compute/virtualmachines" {
			machine.IsArc = false
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
