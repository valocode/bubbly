
data "t1" {
    fields = {
        "f1": "f1val"
    }
}


// These two should be the same
data "t2" {
    joins = ["t1"]
}
data "t2" {
    joins = ["t1"]
}


// These two should be the same
data "t2" {
    fields = {
        "f1": "f1val"
    }
}
data "t2" {
    fields = {
        "f1": "f1val"
    }
}
