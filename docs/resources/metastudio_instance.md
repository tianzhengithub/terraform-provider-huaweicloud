# HuaweiCloud Meta Studio Order Resource

Manages a Meta Studio public cloud service order within HuaweiCloud.

## Example Usage
```hcl
resource "huaweicloud_metastudio_instance" "test" {
  period_type = 2
  period_num = 1
  is_auto_renew = 1
  resource_spec_code = "hws.resource.type.metastudio.avatarmodeling.number"
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

* `is_auto_renew` - (Optional, Int, ForceNew) Specifies whether to auto-renew the vault when it expires.
  Changing this creates a new resource.  
  Valid values are:
  + `0` - Do not renew automatically (default)
  + `1` - Renew automatically

* `resource_spec_code` - (Required, String, ForceNew) Specifies the resource specification code.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `resource_id` - The resource ID of the Meta Studio resource.

* `resource_type` - The type of the resource.

* `business_type` - The business type of the resource.

* `order_id` - The order ID associated with the resource.

* `resource_expire_time` - The expiration time of the resource.

* `status` - The status of the resource.

* `charging_mode` - The charging mode of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes (includes order processing time).
* `delete` - Default is 10 minutes.

## Import

The Meta Studio resource can be imported using `id`, e.g.

## Important Notes

1. **Charging Mode**:
  - For **prePaid** resources (periodic billing), deleting the resource will unsubscribe the service
  - For **ONE_TIME** resources , deleting only removes the resource from Terraform state

2. **Resource Status**:
  - Terraform will automatically wait for the order to complete processing after creation
  - Resource status can be checked through the `status` attribute

3. **Unsupported Operations**:
  - This resource does not support updates after creation
  - Changing any parameter requires recreating the resource
