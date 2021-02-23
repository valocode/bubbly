package store

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/events"
)

// A trigger is a wrapper around a visitFn, which defines what should happen
// at each node of the data tree.
// A trigger can be seen as a mechanism for introducing store-level
// "cause-and-effect", in that it will typically evaluate changes
// to subsections of the store (cause) and generate core.DataBlocks
// representing resultant desired changes to other subsections (effect).
type trigger struct {
	id          string
	description string
	visitFn     visitFn
}

// HandleTriggers evaluates changes to the dataTree against a set of triggers.
// Each trigger contains a visitFn that should "trigger" the injection of
// new dataBlocks into the store on certain conditions.
// Once all triggers have been evaluated, the dataBlocks are converted
// to a dataTree and returned to be saved to the store
func HandleTriggers(tree dataTree, triggers []*trigger) (dataTree, error) {
	triggerBlocks := core.DataBlocks{}
	for _, t := range triggers {
		blocks, err := tree.traverse(t.visitFn)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to use trigger to traverse data tree: %w",
				err)
		}
		triggerBlocks = append(triggerBlocks, blocks...)
	}

	// if the trigger has generated dataBlocks,
	// they should be returned as a dataTree for saving to the store
	if len(triggerBlocks) > 0 {
		triggerTree, err := createDataTree(triggerBlocks)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to create tree of data blocks while handling triggers: %w", err)
		}
		return triggerTree, nil
	}

	return nil, nil
}

var internalTriggers = []*trigger{eventStoreTrigger}

var eventStoreTrigger = &trigger{
	id:          "default/trigger/event_store_trigger",
	description: "update event store upon new/updated entry to resource store",
	visitFn: func(node *dataNode, blocks *core.DataBlocks) error {
		switch node.Data.TableName {
		// if the table is _resource, then we know a resource has been loaded to the store
		case core.ResourceTableName:
			// if the _resource is provided only to provide a FK for another
			// _event entry (and therefore the only field provided is the "id"),
			// then we do not want to update the event store as that will take
			// place _anyway_.
			if len(node.Data.Fields) != 1 && node.Data.IsValidResource() {
				fields := node.Data.Fields

				// get the id of the resource
				id := fields["id"]

				if id.IsNull() {
					return fmt.Errorf(`DataBlock missing required "%s" field`, "id")
				}

				newEventBlocks := core.DataBlocks{
					core.Data{
						TableName: core.ResourceTableName,
						Fields: core.DataFields{
							"id": id,
						},
					},
					core.Data{
						TableName: core.EventTableName,
						Fields: map[string]cty.Value{
							"status": cty.StringVal(events.ResourceCreatedUpdated.String()),
							"time":   cty.StringVal(events.TimeNow()),
						},
						Joins: []string{core.ResourceTableName},
					},
				}

				*blocks = append(*blocks, newEventBlocks...)
			}

		case core.SchemaTableName:
			// TODO
		default:
			// fmt.Printf("trigger has no action for table: %s\n", node.Data.TableName)
		}
		return nil
	},
}
