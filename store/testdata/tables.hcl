// This file contains some tables that are parsed into []core.Table 
// and used as test input for the store.
// The reason behind this is to make the test data more maintainable rather
// than creating the data manually in Go (which ends up creating a lot of bloat)

// Original test for GraphQL resolver's ability 
// to deal with many-to-many relationships.
table "root" {
    field "name" {
        unique = true
        type = string
    }

    table "child_a" {
        field "name" {
            unique = true
            type = string
        }
        table "grandchild_a" {
            field "name" {
                unique = true
                type = string
            }
        }
    }
    table "child_b" {
        // Table root has only one child_b (one-to-one)
        unique = true
        field "name" {
            unique = true
            type = string
        }
    }
    // Regression test: sibling child nodes
    table "child_c" {
        field "name" {
            unique = true
            type = string
        }
    }
}

table "subroot" {
    // subroot belongs to root
    join "root" { }
    field "name" {
        type = string
    }

    table "subroot_a" {
        field "name" {
            unique = true
            type = string
        }
    }
}