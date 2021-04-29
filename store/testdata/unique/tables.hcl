table "t1" {
    field "f1" {
        type = string
        unique = true
    }
    field "f2" {
        type = number
        unique = true
    }
}

table "t2" {
    field "f1" {
        type = string
        unique = true
    }
    join "t1" { unique = true }
}

table "t3" {
    field "f1" {
        type = string
        unique = true
    }
    join "t1" { unique = true }
}

table "t4" {
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}
table "t5" {
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}
table "t6" {
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}

table "t7" {
    field "f1" {
        type = string
        unique = true
    }
    field "f2" {
        type = string
        unique = true
    }
    field "f3" {
        type = string
        unique = true
    }
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}
table "t8" {
    field "f1" {
        type = string
        unique = true
    }
    field "f2" {
        type = string
        unique = true
    }
    field "f3" {
        type = string
        unique = true
    }
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}
table "t9" {
    field "f1" {
        type = string
        unique = true
    }
    field "f2" {
        // TODO: try different types but need fix for this first
        // type = number
        type = string
        unique = true
    }
    field "f3" {
        // TODO: try different types but need fix for this first
        // type = bool
        type = string
        unique = true
    }
    join "t1" { unique = true }
    join "t2" { unique = true }
    join "t3" { unique = true }
}
