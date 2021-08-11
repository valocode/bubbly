
data "t1" {
    fields {
        f1 = "abc"
    }
    data "t2" {
        fields {
            f1 = "abc"
        }
        data "t3" {
            fields {
                f1 = "abc"
            }
        }
    }
    data "t2" {
        fields {
            f1 = "def"
        }
        // No t3
    }
}