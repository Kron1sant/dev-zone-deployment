package yacloudclient

import (
	"context"
	"devZoneDeployment/config"
	"log"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

var sdk = newSdk()

func newSdk() *ycsdk.SDK {
	ctx := context.Background()
	token := config.SecConfig.GetYaCloudToken()
	if token == "" {
		log.Println("yandex token is not specified")
	}
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.OAuthToken(token),
	})
	if err != nil {
		log.Fatalf("failed creating Yandex Cloud SDK: %s", err)
	}

	return sdk
}

func ListInstances() []*compute.Instance {
	folderID := config.AppConfig.GetString("yacloud.folderid")
	if folderID == "" {
		log.Fatal("yandex cloud folder is not specified")
	}

	request := &compute.ListInstancesRequest{
		FolderId: folderID,
	}
	ctx := context.Background()
	op, err := sdk.Compute().Instance().List(ctx, request)
	if err != nil {
		log.Fatalf("failed receiving list of instances: %s", err)
	}
	return op.Instances
}

func InstanseStatus(s compute.Instance_Status) string {
	return compute.Instance_Status_name[int32(s)]
}
