table "t1" {
    field "f1" {
        type = string
    }
    table "t2" {
        field "f1" {
            type = string
        }
        table "t3" {
            field "f1" {
                type = string
            }
        }
    }
}