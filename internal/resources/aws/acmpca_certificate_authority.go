package aws

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"

	"github.com/infracost/infracost/internal/usage"

	"github.com/shopspring/decimal"
)

type AcmpcaCertificateAuthority struct {
	Address         *string
	Region          *string
	MonthlyRequests *int64 `infracost_usage:"monthly_requests"`
}

var AcmpcaCertificateAuthorityUsageSchema = []*schema.UsageItem{{Key: "monthly_requests", ValueType: schema.Int64, DefaultValue: 0}}

func (r *AcmpcaCertificateAuthority) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *AcmpcaCertificateAuthority) BuildResource() *schema.Resource {
	region := *r.Region

	costComponents := []*schema.CostComponent{
		{
			Name:            "Private certificate authority",
			Unit:            "months",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AWSCertificateManager"),
				ProductFamily: strPtr("AWS Certificate Manager"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "usagetype", ValueRegex: strPtr("/PaidPrivateCA/")},
				},
			},
		},
	}

	certificateTierLimits := []int{1000, 9000}
	if r.MonthlyRequests != nil {
		monthlyCertificatesRequests := decimal.NewFromInt(*r.MonthlyRequests)
		certificateTiers := usage.CalculateTierBuckets(monthlyCertificatesRequests, certificateTierLimits)

		if certificateTiers[0].GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (first 1K)", "0", &certificateTiers[0]))
		}

		if certificateTiers[1].GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (next 9K)", "1000", &certificateTiers[1]))
		}

		if certificateTiers[2].GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (over 10K)", "10000", &certificateTiers[2]))
		}
	} else {
		var unknown *decimal.Decimal
		costComponents = append(costComponents, certificateCostComponent(region, "Certificates (first 1K)", "0", unknown))
	}

	return &schema.Resource{
		Name:           *r.Address,
		CostComponents: costComponents, UsageSchema: AcmpcaCertificateAuthorityUsageSchema,
	}
}

func certificateCostComponent(region string, displayName string, usageTier string, monthlyQuantity *decimal.Decimal) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            displayName,
		Unit:            "requests",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: monthlyQuantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("aws"),
			Region:        strPtr(region),
			Service:       strPtr("AWSCertificateManager"),
			ProductFamily: strPtr("AWS Certificate Manager"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "usagetype", ValueRegex: strPtr("/PrivateCertificatesIssued/")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			StartUsageAmount: strPtr(usageTier),
		},
	}
}
