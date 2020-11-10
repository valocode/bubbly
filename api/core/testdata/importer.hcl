
//  defines a data source and a format for the data, and returns a cty.Value
// representation of the external data.
resource "importer" "sonarqube" {
    api_version = "v1"
    spec {
        input "file" {}
        type = "json"
        source {
            file = self.input.file
            format = object({
                issues: list(object({
                    engineId: string,
                    ruleId: string,
                    severity: string,
                    type: string,
                    primaryLocation: object({
                        message: string,
                        filePath: string,
                        textRange: object({
                            startLine: number,
                            endLine: number,
                            startColumn: number,
                            endColumn: number
                        })
                    })
                }))
            })
        }
    }
}
