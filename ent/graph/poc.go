package graph

import "github.com/valocode/bubbly/ent"

type HCLGraph struct {
	CodeScanNode *CodeScanNode `json:"code_scan" hcl:"code_scan,block"`
}

// type Graph struct {
// 	// Nodes         []Node          `json:"node"` // this would require custom marshalling and unmarshalling...bore!
// 	CodeScanNodes  []CodeScanNode  `json:"code_scan"`
// 	CodeIssueNodes []CodeIssueNode `json:"code_issue"`
// }

// type Node struct {
// 	Operation string
// 	ID        *int
// }

// type CodeScanNode struct {
// 	// NodeOperation string
// 	Node
// 	Tool   *string         `json:"tool" hcl:"tool,attr"`
// 	Issues []CodeIssueNode `json:"issues" hcl:"issue,block"`
// }

func (cs *CodeScanNode) Process(client *ent.Client) error {
	// Must handle required edges...

	switch cs.Node.Operation {
	// create:
	case "":
		create := client.CodeScan.Create()
		if cs.Tool != nil {
			create.SetTool(*cs.Tool)
		}
		// delete:
		// update:
		// query:
		// error:
	}
	return nil
}

// type CodeIssueNode struct {
// 	Rule    *string
// 	Message *string
// }
