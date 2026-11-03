# Top-level namespaces of a warehouse.
data "lakekeeper_namespaces" "roots" {
  warehouse_id = "a8a6b1b8-1c0e-4c3b-9f0a-2e5d7f3c9b11"
}

# Direct children of a nested namespace.
data "lakekeeper_namespaces" "instruments" {
  warehouse_id = "a8a6b1b8-1c0e-4c3b-9f0a-2e5d7f3c9b11"
  parent       = "reference.instruments"
}
