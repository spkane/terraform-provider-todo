---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "todo Data Source - terraform-provider-todo"
subcategory: ""
description: |-
  The todo data source allows you to retrieve information about a particular todo/reminder.
---

# todo (Data Source)

The todo data source allows you to retrieve information about a particular todo/reminder.

```terraform
data "todo" "not_mine" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `completed` (Boolean)
- `description` (String)
- `id` (Number) The ID of this resource.
