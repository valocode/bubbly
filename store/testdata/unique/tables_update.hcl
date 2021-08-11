table "tu1" {
    field "f1" {
        type = string
    }
}

table "tu2" {
    field "f1" {
        type = string
        unique = true
    }
    join "tu1" {}
}