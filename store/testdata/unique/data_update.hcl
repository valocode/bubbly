data "tu2" {
    fields {
        f1 = "unique"
    }
}

data "tu1" {
    fields {
        _id = self.data.tu2.tu1_id
        f1 = "not_unique"
    }
}

data "tu2" {
    fields {
        _id = self.data.tu2._id
    }
    joins = ["tu1"]
    policy = "update"
}
