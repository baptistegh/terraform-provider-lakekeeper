data "lakekeeper_warehouses" "this" {
  project_id = "9d25e96e-601a-44ea-badc-74b8c896b4f0"
}

locals {
  active_warehouse_ids = { for w in data.lakekeeper_warehouses.this.warehouses : w.name => w.id if w.active }
}
