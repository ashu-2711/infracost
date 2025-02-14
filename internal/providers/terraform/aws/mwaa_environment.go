package aws

import (
	"github.com/infracost/infracost/internal/resources/aws"
	"github.com/infracost/infracost/internal/schema"
	"github.com/tidwall/gjson"
	"strings"
)

func GetMWAAEnvironmentRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "aws_mwaa_environment",
		RFunc: NewMWAAEnvironment,
	}
}

func NewMWAAEnvironment(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()

	size := "mw1.small"
	if d.Get("environment_class").Type != gjson.Null {
		size = d.Get("environment_class").String()
	}

	size = strings.ToLower(size)
	size = strings.ReplaceAll(size, "mw1.", "")
	size = strings.Title(size)

	a := &aws.MWAAEnvironment{
		Address: d.Address,
		Region:  region,
		Size:    size,
	}
	a.PopulateUsage(u)

	return a.BuildResource()
}
