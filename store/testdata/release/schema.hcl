
// release is a way to logically mark a "version" of some release item(s)
table "release" {
    field "name" {
        type = string
        unique = true
    }
    field "version" {
        type = string
        unique = true
    }

    join "project" { unique = true }
}

// release_item is used to represent what we are releasing in a single release
// and can be of different types: git, artifact or release.
// Based on the type it should have a join to one of those tables
table "release_item" {

    // type should be one of git (commit), artifact or release
    field "type" {
        type = string
    }

    // Join to release. A release can have one or more release_items.
    // A release_item can belong to only one release, because it can have
    // crtieria associated with it, which are specific to a release
    join "release" {
        unique = true
    }

    // Join to the different item tables with a one-to-one relationship.
    // Only at most and at least one of these joins should exist, based on the
    // "type" field.
    // This means each commit/artifact can have only one release_item
    join "commit" {
        single = true
        unique = true
    }
    join "artifact" {
        single = true
        unique = true
    }
    // TODO: this is a problem because this creates a second join from release_item
    // to release... It could be solved by adding an alias
    // join "release" {
    //     alias = "item_release"
    //     single = true
    //     unique = true
    // }
}

// release_entry is used to record/log an event performed on a release,
// such as running of unit tests, or the creation of an artifact.
// release_entry is created by running a criteria and should contain the output
// from the running of that event
table "release_entry" {
    field "name" {
        type = string
        unique = true
    }
    field "result" {
        type = bool
    }

    // TODO: what other fields do we want to store? Probably something saying
    // *why* the criteria failed (a reason) and also perhaps the GraphQL
    // query used so that we could fetch the data? E.g.
    // field "query" { type = string}
    // field "reason" { type = string}

    // A release_entry is unique by it's join to a release
    join "release" { unique = true }
    // Join on the _resource criteria that created this entry release_entry
    join "_resource" { }
}

table "release_stage" {
    field "name" {
        type = string
        unique = true
    }

    join "release" { 
        unique = true
    }
}

table "release_criteria" {
    field "entry_name" {
        type = string
        unique = true
    }
    join "release_stage" { unique = true }
}


// #############################
// Below tables should already exist elsewhere for the release_item types
// #############################

// _resource table already exists in bubbly, so no need to duplicate it, but we
// need it for testing this schema
table "_resource" {
    field "id" {
        type = string
        unique = true
    }
}

table "project" {
    field "name" {
        type = string
        unique = true
    }
}

table "repo" {
    field "name" {
        type = string
        unique = true
    }

    table "branch" {
        unique = true
        field "name" {
            type = string
            unique = true
        }

        // A specific commit/version in a git repository
        table "commit" {
            field "id" {
                type = string
                unique = true
            }
            field "tag" {
                type = string
            }
            // Would be really cool to store the time of a commit, and then we can
            // track how long it takes to do things, e.g. time to deploy
            field "time" {
                type = string
            }
        }
    }


    join "project" {}
    // join "branch" {
    //     alias = "default_branch"
    // }
}

table "artifact" {
    // The sha256 of an artifact shall uniquely identify it
    field "sha256" {
        type = string
        unique = true
    }
    // A url, or path to a docker image. Should always start with a type, e.g.
    // https:// or docker://
    field "location" {
        type = string
    }
}

table "test_run" {
    field "name" {
        type = string
    }
    field "result" {
        type = bool
    }

    table "test_case" {
        field "name" {
            type = string
        }
    }

    join "release" {}
}
