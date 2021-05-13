// #############################
// Let's create the dummy data that needs to exist first
// #############################

data "project" {
    fields = {
        "name": "bubbly"
    }
}

data "repo" {
    fields = {
        "name": "github.com/valocode/bubbly"
    }
    data "branch" {
        fields = {
            "name": "main"
        }
        data "commit" {
            fields = {
                "id": "asdasdasdasdasd",
                "tag": "0.1.23"
            }
        }
    }
}

data "artifact" {
    fields = {
        "sha256": "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
        "location": "docker://valocode/bubbly:0.1.23"
    }
}

data "_resource" {
    fields = {
        "id": "criteria/unit_tests"
    }
}

// #############################
// Let's create an example release
// #############################

data "project" {
    fields = {
        "name": "bubbly"
    }
    policy = "reference"
}

// Create a reference to a git repo commit that we can join to
data "commit" {
    fields = {
        "id": "asdasdasdasdasd"
    }
    // TODO: we should implement something so that it does not actually create
    // the data block if it doesn't exist, but instead errors if it doesn't exist
    // or then selects the required fields for us to join against it
    policy = "reference"
}

data "release" {
    fields = {
        "name": "bubbly"
    }
    data "release_item" {
        fields = {
            "type": "git"
        }
        joins = ["commit"]
    }
    data "release_stage" {
        fields = {
            "name": "Testing"
        }
        data "release_criteria" {
            fields = {
                "entry_name": "unit_test"
            }
        }
    }
    joins = ["project"]
}

// #############################
// Let's create an example release_entry
// #############################

// In order to do that, we need to fetch the release_item so that we can make a join.
// Remember that creating the release_entry will happen in CI or somewhere and not
// in the same command that creates the release...
// A release_item is uniquely identified by the joins (commit_id, artifact_id, release_id)
// so we need to also fetch the data corresponding to the join table
data "commit" {
    fields = {
        "id": "asdasdasdasdasd",
    }
    policy = "reference"
}

data "release_item" {
    joins = ["commit", "release"]    
    policy = "reference"
}

data "release_entry" {
    fields = {
        "name": "unit_test",
        "result": true
    }

    // Join to the release_item we reference above
    joins = ["release"]
}


// #############################
// Let's load a testrun
// #############################

data "project" {
    fields = {
        "name": "bubbly"
    }
    policy = "reference"
}

data "release" {
    fields = {
        "name": "bubbly"
    }
    joins = ["project"]
    policy = "reference"
}

data "testrun" {
    fields = {
        "name": "unit tests",
        "result": true
    }
    data "test_case" {
        fields = { "name": "test_case_1" }
    }
    data "test_case" {
        fields = { "name": "test_case_2" }
    }
    joins = ["release"]
}
