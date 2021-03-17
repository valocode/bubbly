// This file contains some data blocks that are parsed into []core.Data
// and used as test input for the store.
// The reason behind this is to make the test data more maintainable rather
// than creating the data manually in Go (which ends up creating a lot of bloat)

data "root" {
    fields = {
        "name": "first_root"
    }
}

data "child_a" {
    joins = ["root"]
    fields = {
        "name": "first_child"
    }
    data "grandchild_a" {
        fields = {
            "name": "first_grandchild"
        }
    }
}

// Regression test: sibling child nodes
data "child_c" {
    joins = ["root"]
    fields = {
        "name": "sibling_child"
    }
}

data "grandchild_a" {
    joins = ["child_a"]
    fields = {
        "name": "second_grandchild"
    }
}

data "root" {
    fields = {
        "name": "second_root"
    }
}