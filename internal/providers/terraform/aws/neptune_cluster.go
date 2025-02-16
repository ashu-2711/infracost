package aws

import (
	"github.com/infracost/infracost/internal/resources/aws"
	"github.com/infracost/infracost/internal/schema"
)

func getNeptuneClusterRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "aws_neptune_cluster",
		RFunc: NewNeptuneCluster,
	}
}
func NewNeptuneCluster(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	r := &aws.NeptuneCluster{Address: strPtr(d.Address), Region: strPtr(d.Get("region").String())}
	if !d.IsEmpty("backup_retention_period") {
		r.BackupRetentionPeriod = intPtr(d.Get("backup_retention_period").Int())
	}
	r.PopulateUsage(u)
	return r.BuildResource()
}
