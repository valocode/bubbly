package store

//
// import (
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"os"
// 	"reflect"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/verifa/bubbly/api"
// 	v1 "github.com/verifa/bubbly/api/v1"
// 	"github.com/verifa/bubbly/env"
// 	"github.com/verifa/bubbly/parser"
//
// 	"github.com/go-pg/pg/v10"
// 	"github.com/verifa/bubbly/api/core"
// )
//
// func Test_postgres_PutResource(t *testing.T) {
// 	bCtx := env.NewBubblyContext()
// 	resParser := api.NewParserType()
// 	err := parser.ParseFilename(bCtx, "../integration/testdata/sonarqube/sonarqube-pipeline.bubbly", resParser)
// 	assert.NoError(t, err)
//
// 	err = resParser.CreateResources(bCtx)
// 	assert.NoError(t, err)
// 	var res core.Resource
// 	var resJson []byte
// 	var pipelineRun *v1.PipelineRun
// 	for _, r := range resParser.Resources {
// 		if r.Kind() == core.PipelineRunResourceKind {
// 			resJson, _ = json.Marshal(r)
// 			pipelineRun = r.(*v1.PipelineRun)
// 			res = r
// 		}
// 	}
// 	if res == nil || pipelineRun == nil {
// 		assert.Fail(t, "no pipeline_run was found in bubbly file")
// 		return
// 	}
//
// 	db := pg.Connect(&pg.Options{
// 		Addr:     getenv("postgres", "localhost") + ":5432", // TODO set
// 		User:     "postgres",
// 		Password: "postgres",
// 		Database: "bubbly",
// 	})
//
// 	type fields struct {
// 		db    *pg.DB
// 		types map[string]schemaType
// 		ctx   context.Context
// 	}
// 	type args struct {
// 		id       string
// 		val      string
// 		resBlock core.ResourceBlock
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Insert Pipeline_Run into DB",
// 			args: args{
// 				id:       res.String(),
// 				val:      string(resJson),
// 				resBlock: *pipelineRun.ResourceBlock,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &postgres{
// 				db:    db,
// 				types: tt.fields.types,
// 				ctx:   tt.fields.ctx,
// 			}
// 			if err := p.PutResource(tt.args.id, tt.args.val, tt.args.resBlock); (err != nil) != tt.wantErr {
// 				t.Errorf("PutResource() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
// func Test_postgres_GetResourcesByKind(t *testing.T) {
// 	db := pg.Connect(&pg.Options{
// 		Addr:     getenv("postgres", "localhost") + ":5432",
// 		User:     "postgres",
// 		Password: "postgres",
// 		Database: "bubbly",
// 	})
//
// 	type fields struct {
// 		db    *pg.DB
// 		types map[string]schemaType
// 		ctx   context.Context
// 	}
// 	type args struct {
// 		resourceKind core.ResourceKind
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    io.Reader
// 		wantErr bool
// 	}{
// 		{
// 			name: "Select pipeline_runs",
// 			fields: fields{
// 				db: db,
// 			},
// 			args: args{
// 				resourceKind: core.PipelineRunResourceKind,
// 			},
// 			wantErr: false,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &postgres{
// 				db:    tt.fields.db,
// 				types: tt.fields.types,
// 				ctx:   tt.fields.ctx,
// 			}
// 			got, err := p.GetResourcesByKind(tt.args.resourceKind)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetResourcesByKind() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetResourcesByKind() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func getenv(key, fallback string) string {
// 	value := os.Getenv(key)
// 	if len(value) == 0 {
// 		return fallback
// 	}
// 	return value
// }
