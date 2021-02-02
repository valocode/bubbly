// This file contains some data blocks that are parsed into []core.Data
// and used as test input for the store.
// The reason behind this is to make the test data more maintainable rather
// than creating the data manually in Go (which ends up creating a lot of bloat)

data "root" {
    fields = {
        "name": "test_value"
    }
}

data "child_a" {
    joins = ["root"]
    fields = {
        "name": "test_value"
    }
    data "grandchild_a" {
        fields = {
            "name": "nested_value"
        }
    }
}

data "grandchild_a" {
    joins = ["child_a"]
    fields = {
        "name": "join_value"
    }
}
