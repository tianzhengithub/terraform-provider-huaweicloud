package metastudio

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"strings"
)

func ResourceMetaStudio() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetaStudioCreate,
		ReadContext:   resourceMetaStudioRead,
		UpdateContext: resourceMetaStudioUpdate,
		DeleteContext: resourceMetaStudioDelete,
		Schema: map[string]*schema.Schema{
			"cloud_services": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 0,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_auto_pay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 1),
							Default:      0,
							Description:  `Specifies whether to pay automatically.`,
						},
						"period_type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{2, 3, 6}),
							Description:  `Specifies the charging period unit`,
						},
						"period_num": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 2147483647),
							Description:  `Specifies the number of periods to purchase.`,
						},
						"is_auto_renew": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 1),
							Default:      0,
							Description:  `Specifies whether to auto-renew the vault when it expires.`,
						},
						"subscription_num": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 2147483647),
							Description:  `Specifies the number of subscription`,
						},
						"resource_spec_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the resource specification`,
						},
					},
				},
				Description: "Array of PublicCloudServiceOrder objects",
			},
		},
	}
}

func resourceMetaStudioCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/mss/public/orders"
	)
	client, err := cfg.MetaStudioClient(region)
	if err != nil {
		return diag.Errorf("error creating MetaStudio client: %s", err)
	}
	orderId, err := resourceCreate(client, d, httpUrl)
	if err != nil {
		return diag.Errorf("error creating MetaStudio: %s", err)
	}
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	timeout := d.Timeout(schema.TimeoutCreate)
	// wait for order complete
	if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeout); err != nil {
		return diag.FromErr(err)
	}
	return nil

}

func resourceCreate(client *golangsdk.ServiceClient, d *schema.ResourceData, httpUrl string) (string, error) {
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"cloud_services": d.Get("cloud_services").([]interface{}),
		},
	}
	resp, err := client.Request("POST", requestPath, &createOpt)
	if err != nil {
		return "", err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}
	orderId := utils.PathSearch("order_id", respBody, "").(string)
	return orderId, nil
}

func resourceMetaStudioDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceMetaStudioUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceMetaStudioRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
