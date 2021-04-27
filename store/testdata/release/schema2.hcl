table "t1" {
    field "f1" {
        type = string
        unique = true
    }

    table "t2" {
        // Make the join part of the unique constraint
        unique = true
        field "f1" {
            type = string
            unique = true
        }
    }
}