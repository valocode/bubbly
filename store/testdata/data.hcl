// This file contains some data blocks that are parsed into []core.Data
// and used as test input for the store.
// The reason behind this is to make the test data more maintainable rather
// than creating the data manually in Go (which ends up creating a lot of bloat)

data "root" {
    field "name" {
        value = "test_value"
    }
}

data "child_a" {
    join "root" {
        value = self.data.root._id
    }
    field "name" {
        value = "test_value"
    }
    data "grandchild_a" {
        field "name" {
            value = "nested_value"
        }
    }
}

data "grandchild_a" {
    join "child_a" {
        value = self.data.child_a._id
    }
    field "name" {
        value = "join_value"
    }
}