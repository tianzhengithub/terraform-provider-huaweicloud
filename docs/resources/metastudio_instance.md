
# HuaweiCloud Meta Studio Resource

Manages a Meta Studio resource within HuaweiCloud, supporting pre-paid billing mode with auto-renewal capabilities.

## Example Usage
```hcl
resource "huaweicloud_metastudio_instance" "test" {
  period_type = 2
  period_num = 1
  is_auto_renew = 0
  resource_spec_code = "hws.resource.type.metastudio.modeling.avatarlive.channel"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `period_type` - (Required, Int, ForceNew) Specifies the charging period unit.
  Changing this creates a new resource.  
  Valid values are:
  + `2` - Month
  + `3` - Year
  + `6` - Week

* `period_num` - (Required, Int, ForceNew) Specifies the number of periods to purchase.
  Changing this creates a new resource.  
  Value range: `1` to `2147483647`.

* `is_auto_renew` - (Optional, Int, ForceNew) Specifies whether to auto-renew the resource when it expires.
  Changing this creates a new resource.  
  Valid values are:
  + `0` - Do not renew automatically (default)
  + `1` - Renew automatically

* `resource_spec_code` - (Required, String, ForceNew) Specifies the resource specification code for user-purchased cloud service products. For details, see [Resource Types](https://support.huaweicloud.com/api-metastudio/metastudio_02_0042.html).
  Changing this creates a new resource.

* `enable_force_new` - (Optional, String) Internal parameter (not recommended for user configuration).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `order_id` - The order ID associated with the resource.

* `resource_expire_time` - The expiration time of the resource.

* `business_type` - The business type of the resource.

* `sub_resource_type` - The sub-resource type.

* `is_sub_resource` - Indicates whether it is a sub-resource.

* `charging_mode` - The billing mode of the resource (e.g., `PERIODIC` for pre-paid).

* `amount` - The total amount of the resource.

* `usage` - The usage amount of the resource.

* `status` - The status of the resource:
  + `0` - Normal
  + `1` - Frozen

* `unit` - The unit of measurement for the resource amount.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes (includes order processing time).
* `delete` - Default is 10 minutes.

## Import

The Meta Studio resource can be imported using `id`, e.g.

## Important Notes

1. **Auto-Pay Agreement**:
  - Terraform will automatically sign the auto-pay agreement during resource creation
  - No user intervention is required for this process

2. **Billing Modes**:
  - **Periodic (pre-paid)**: Deleting the resource will unsubscribe the service
  - **Other modes**: Deleting only removes the resource from Terraform state

3. **Resource Status**:
  - Terraform automatically waits for order completion after creation
  - Resource status can be monitored through the `status` attribute

4. **Unsupported Operations**:
  - This resource does not support updates after creation
  - Changing any parameter requires recreating the resource

5. **Deletion Behavior**:
  - Pre-paid resources are unsubscribed upon deletion
  - Non-prepaid resources are only removed from state management