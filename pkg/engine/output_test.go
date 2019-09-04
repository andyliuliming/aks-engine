// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/Azure/aks-engine/pkg/api"
	"github.com/Azure/aks-engine/pkg/helpers"
	"github.com/Azure/aks-engine/pkg/i18n"
)

func TestWriteTLSArtifacts(t *testing.T) {

	cs := api.CreateMockContainerService("testcluster", "1.7.12", 1, 2, true)
	writer := &ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}
	dir := "_testoutputdir"
	defaultDir := fmt.Sprintf("%s-%s", cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.GetClusterID())
	defaultDir = path.Join("_output", defaultDir)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(defaultDir)

	// Generate apimodel and azure deploy artifacts without certs
	err := writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", dir, false, false)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	expectedFiles := []string{"apimodel.json", "azuredeploy.json", "azuredeploy.parameters.json"}

	for _, f := range expectedFiles {
		if _, err = os.Stat(dir + "/" + f); os.IsNotExist(err) {
			t.Fatalf("expected file %s/%s to be generated by WriteTLSArtifacts", dir, f)
		}
	}

	os.RemoveAll(dir)

	// Generate parameters only and certs
	err = writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", "", true, true)
	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	if _, err = os.Stat(defaultDir + "/apimodel.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/apimodel.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	if _, err = os.Stat(defaultDir + "/azuredeploy.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	expectedFiles = []string{"azuredeploy.parameters.json", "ca.crt", "ca.key", "apiserver.crt", "apiserver.key", "client.crt", "client.key", "etcdclient.key", "etcdclient.crt", "etcdserver.crt", "etcdserver.key", "etcdpeer0.crt", "etcdpeer0.key", "kubectlClient.crt", "kubectlClient.key"}

	for _, f := range expectedFiles {
		if _, err = os.Stat(defaultDir + "/" + f); os.IsNotExist(err) {
			t.Fatalf("expected file %s/%s to be generated by WriteTLSArtifacts", dir, f)
		}
	}

	kubeDir := path.Join(defaultDir, "kubeconfig")
	if _, err = os.Stat(kubeDir + "/" + "kubeconfig.eastus.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/kubeconfig/kubeconfig.eastus.json to be generated by WriteTLSArtifacts", defaultDir)
	}
	os.RemoveAll(defaultDir)

	// Generate files with custom cloud profile in configuration
	csCustom := api.CreateMockContainerService("testcluster", "1.11.6", 1, 2, false)
	csCustom.Location = "customlocation"
	csCustom.Properties.CustomCloudProfile = &api.CustomCloudProfile{
		Environment: &azure.Environment{
			Name:                         "azurestackcloud",
			ManagementPortalURL:          "managementPortalURL",
			PublishSettingsURL:           "publishSettingsURL",
			ServiceManagementEndpoint:    "serviceManagementEndpoint",
			ResourceManagerEndpoint:      "resourceManagerEndpoint",
			ActiveDirectoryEndpoint:      "activeDirectoryEndpoint",
			GalleryEndpoint:              "galleryEndpoint",
			KeyVaultEndpoint:             "keyVaultEndpoint",
			GraphEndpoint:                "graphEndpoint",
			ServiceBusEndpoint:           "serviceBusEndpoint",
			BatchManagementEndpoint:      "batchManagementEndpoint",
			StorageEndpointSuffix:        "storageEndpointSuffix",
			SQLDatabaseDNSSuffix:         "sqlDatabaseDNSSuffix",
			TrafficManagerDNSSuffix:      "trafficManagerDNSSuffix",
			KeyVaultDNSSuffix:            "keyVaultDNSSuffix",
			ServiceBusEndpointSuffix:     "serviceBusEndpointSuffix",
			ServiceManagementVMDNSSuffix: "serviceManagementVMDNSSuffix",
			ResourceManagerVMDNSSuffix:   "resourceManagerVMDNSSuffix",
			ContainerRegistryDNSSuffix:   "containerRegistryDNSSuffix",
			TokenAudience:                "tokenAudience",
		},
	}
	csCustom.SetPropertiesDefaults(PropertiesDefaultsOptions{
		IsScale:    false,
		IsUpgrade:  false,
		PkiKeySize: helpers.DefaultPkiKeySize,
	})
	err = writer.WriteTLSArtifacts(csCustom, "vlabs", "fake template", "fake parameters", "", true, false)
	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	expectedFiles = []string{"apimodel.json", "azuredeploy.json", "azuredeploy.parameters.json", "ca.crt", "ca.key", "apiserver.crt", "apiserver.key", "client.crt", "client.key", "etcdclient.key", "etcdclient.crt", "etcdserver.crt", "etcdserver.key", "etcdpeer0.crt", "etcdpeer0.key", "kubectlClient.crt", "kubectlClient.key"}

	for _, f := range expectedFiles {
		if _, err = os.Stat(defaultDir + "/" + f); os.IsNotExist(err) {
			t.Fatalf("expected file %s/%s to be generated by WriteTLSArtifacts", dir, f)
		}
	}

	kubeDirCustom := path.Join(defaultDir, "kubeconfig")
	if _, err = os.Stat(kubeDirCustom + "/" + "kubeconfig." + csCustom.Location + ".json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/kubeconfig/kubeconfig.%s.json to be generated by WriteTLSArtifacts", csCustom.Location, defaultDir)
	}
	os.RemoveAll(defaultDir)

	// Generate certs with all kubeconfig locations
	cs.Location = ""
	err = writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", "", true, false)
	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	for _, region := range helpers.GetAzureLocations() {
		if _, err := os.Stat(kubeDir + "/" + "kubeconfig." + region + ".json"); os.IsNotExist(err) {
			t.Fatalf("expected kubeconfig for region %s to be generated by WriteTLSArtifacts", region)
		}
	}
}
