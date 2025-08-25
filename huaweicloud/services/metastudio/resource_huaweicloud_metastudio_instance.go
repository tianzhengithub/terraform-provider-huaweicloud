package metastudio

import (
	"context"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"strings"
	"time"
)

func ResourceMetaStudio() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetaStudioCreate,
		ReadContext:   resourceMetaStudioRead,
		UpdateContext: resourceMetaStudioUpdate,
		DeleteContext: resourceMetaStudioDelete,
		Schema: map[string]*schema.Schema{
			"period_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 6}),
				Description:  `Specifies the charging period unit`,
			},
			"period_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
				Description:  `Specifies the number of periods to purchase.`,
			},
			"is_auto_renew": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Default:      0,
				Description:  `Specifies whether to auto-renew the vault when it expires.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `Specifies the resource specification code`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
		return diag.Errorf("the order (%s) is not completed while creating metaStudio : %v", orderId, err)
	}
	resourceId, err := common.WaitOrderAllResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)
	return resourceMetaStudioRead(ctx, d, meta)
}

func resourceCreate(client *golangsdk.ServiceClient, d *schema.ResourceData, httpUrl string) (string, error) {
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateMetaStudioParams(d)),
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

func buildCreateMetaStudioParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_services": []interface{}{
			map[string]interface{}{
				"is_auto_pay":        1,
				"period_type":        d.Get("period_type"),
				"period_num":         d.Get("period_num"),
				"is_auto_renew":      d.Get("is_auto_renew"),
				"subscription_num":   1,
				"resource_spec_code": d.Get("resource_spec_code"),
			},
		},
	}
	return bodyParams
}

func resourceMetaStudioRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.MetaStudioClient(region)
	if err != nil {
		return diag.Errorf("error creating MetaStudio client: %s", err)
	}
	resource, diagResult := GetResourceDetail(client, d)
	if diagResult != nil {
		return diagResult
	}
	mErr := multierror.Append(nil,
		d.Set("resource_id", utils.PathSearch("resource_id", resource, nil).(string)),
		d.Set("resource_type", utils.PathSearch("resource_type", resource, nil).(string)),
		d.Set("business_type", utils.PathSearch("business_type", resource, nil).(string)),
		d.Set("order_id", utils.PathSearch("order_id", resource, "").(string)),
		d.Set("resource_expire_time", utils.PathSearch("resource_expire_time", resource, nil).(string)),
		d.Set("status", utils.PathSearch("status", resource, nil).(int32)),
		d.Set("charging_mode", utils.PathSearch("charging_mode", resource, nil).(string)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting MetaStudio resource fields: %s", mErr)
	}
	return nil
}

func GetResourceDetail(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	httpUrl := "/v1/{project_id}/tenants/resources"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?resource_source=PURCHASED&resource_id=%v", requestPath, d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, diag.Errorf("error querying MetaStudio resource detail: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, diag.Errorf("error querying MetaStudio resource detail: %s", err)
	}
	resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	if len(resources) < 1 {
		resourceID := d.Id()
		d.SetId("")
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s don't exist", resourceID),
			},
		}
	}
	return resources[0], nil
}

func resourceMetaStudioDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.MetaStudioClient(region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}
	resourceId := d.Id()
	if d.Get("charging_mode").(string) == "PERIODIC" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{resourceId}); err != nil {
			return diag.Errorf("error unsubscribing meta studio resource (%s): %s",
				resourceId, err)
		}
		if err := waitingForResourceDeleteCompleted(ctx, client, d); err != nil {
			return diag.Errorf("error waiting for Workspace APP server (%s) deleted: %s", d.Id(), err)
		}
	} else {
		errorMsg := `This resource is a one-time action resource. Deleting this 
resource will not change the current resource status, but will only remove the resource information from the 
tfstate file.`
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorMsg,
			},
		}
	}
	return nil
}

func waitingForResourceDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) interface{} {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resourceDetail, diagResult := GetResourceDetail(client, d)
			if resourceDetail != nil {
				return resourceDetail, "PENDING", nil
			} else if diagResult[0].Summary == "Resource not found" {
				return "deleted", "COMPLETED", nil
			} else {
				return nil, "ERROR", fmt.Errorf(diagResult[0].Summary)
			}
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceMetaStudioUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource doesn't support update operation.
	return nil
}
