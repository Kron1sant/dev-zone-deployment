package sessions

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db"
	"devZoneDeployment/db/mongodb"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_contextSession_UpdateListVirtualMachinesFromCloud(t *testing.T) {
	type fields struct {
		ctx *gin.Context
		ds  db.DataActions
		ui  api.UserIdentity
	}
	testAdmin := api.UserIdentity{
		Id:       9999,
		IsAdmin:  true,
		Username: "TestAdmin9999",
		Empty:    false,
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "try getting VMs from YaCloud and write it too DB",
			fields: fields{
				ctx: nil,
				ds:  mongodb.UseMongoDBSource(),
				ui:  testAdmin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sctx := &contextSession{
				ctx: tt.fields.ctx,
				ds:  tt.fields.ds,
				ui:  tt.fields.ui,
			}
			if err := sctx.UpdateListVirtualMachinesFromCloud(); (err != nil) != tt.wantErr {
				t.Errorf("contextSession.UpdateListVirtualMachinesFromCloud() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
