# Release Readiness Design

```hcl
resource "release" "bubbly" {
    spec {
        input "git" {}
        input "artifact" {}
        // more inputs...

        // There must be at least one of these
        git {
            repo = self.input.git.repo_name
            commit = self.input.git.commit
        }
        /* TODO LATER... think about complexity with multiple stages * items
        artifact {
            sha = self.input.artifact.sha
        }
        release {
            // TODO
        }
        */
        
        stage "Artifact Published" {
            criterion = [
                "artifact_published"
            ]
        }
        stage "Static Code Analysis" {
            criterion = [
                "crtieria/gosec"
            ]
        }
    }

}

```
