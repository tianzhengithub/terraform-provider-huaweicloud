# HuaweiCloud Meta Studio Order Resource

Manages a Meta Studio public cloud service order within HuaweiCloud.

## Example Usage
```hcl
resource "huaweicloud_meta_studio" "order" {
  cloud_services = [
    {
      period_type        = 2
      period_num         = 12
      subscription_num   = 2
      resource_spec_code = "spec.code.standard"
      is_auto_pay        = 1
      is_auto_renew      = 1
    },
    {
      period_type        = 3
      period_num         = 6
      subscription_num   = 1
      resource_spec_code = "spec.code.premium"
      is_auto_pay        = 0
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the order. If omitted, the provider-level region will be used.

* `cloud_services` - (Required, List) Specifies the list of public cloud service orders. Each cloud service configuration must include the parameters below. Minimum of 0 items, maximum of 100 items.

The `cloud_services` block supports:

* `period_type` - (Required, Int) Specifies the charging period unit. Valid values are:  
  `2` - Month  
  `3` - Year  
  `6` - Week

* `period_num` - (Required, Int) Specifies the number of periods to purchase. Value range: 1 to 2147483647.

* `subscription_num` - (Required, Int) Specifies the number of subscriptions. Value range: 1 to 2147483647.

* `resource_spec_code` - (Required, String) Specifies the resource specification code.

* `is_auto_pay` - (Optional, Int) Specifies whether to pay automatically.  
  `0` - Manual payment (default)  
  `1` - Automatic payment

* `is_auto_renew` - (Optional, Int) Specifies whether to auto-renew the service when it expires.  
  `0` - Do not renew automatically (default)  
  `1` - Renew automatically

## Attribute Reference

* `id` - The order ID in UUID format (populated after successful creation).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes (includes order completion wait time).

## Note

* After creating the order, Terraform will automatically wait for the order to complete processing.
* Once created, this resource does not support update or delete operations through Terraform.
