
data "t1" {
    fields = {
        "f1": "f1val"
    }
}

data "t2" {
    joins = ["t1"]
}

data "t3" {
    fields = {
        "f1": "f1val"
    }
}

data "t4" { joins = ["t1"] }
data "t5" { joins = ["t1", "t2"] }
data "t6" { joins = ["t1", "t2", "t3"] }
data "t7" { 
    fields = { "f1": "f1val" }
    joins = ["t1"]
}
data "t8" { 
    fields = { "f2": "2" }
    joins = ["t1", "t2"]
}
data "t9" { 
    fields = {
        "f2": "2"
        "f3": "true"
    }
    joins = ["t1", "t3"]
}
