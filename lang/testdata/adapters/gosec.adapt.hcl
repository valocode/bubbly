
// bubbly log security --tool gosec --json gosec_*.json
// bubbly log security --tool gosec --xml gosec_*.json

adapter {
    name = "gosec"
    // optional: the tag which uniquely identifies this adapter (default: default)
    tag = ""
    type = "json" // xml, http, regexp
    
    // optional: operation configuration based on the type of the adapter
    // operation {
        // optional: override the default file (otherwise gosec.bubbly.json)
        // file = ""
        // optional: returns a string that should be parsed in case of invalid json, for example
        // preprocess = "[${join(",", split("\n", trimspace(raw_data)))}]"
        // optional for json/yaml, otherwise the type is implied
        // structure = object({
        //     Issues: list(object({
        //         severity: string,
        //         cwe: object({
        //             id: string,
        //             url: string,
        //         }),
        //         rule_id: string,
        //         details: string,
        //     })),
        // })
    // }

    mapping {
        code_scan {
            tool = "gosec"
            dynamic "issues" {
                for_each = data.Issues
                iterator = issue
                content {
                    type = "security"
                    rule_id = issue.value.rule_id
                    message = issue.value.details
                    severity = lower(issue.value.severity)
                }
            }
        }
    }
}
