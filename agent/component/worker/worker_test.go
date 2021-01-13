package worker

// import (
// 	"context"
// 	"testing"
//
// 	zerzerolog "github.com/rs/zerolog"
//
// 	"github.com/stretchr/testify/assert"
//
// 	"github.com/verifa/bubbly/env"
// )
//
// func TestWorker_Run(t *testing.T) {
// 	bCtx := env.NewBubblyContext()
// 	bCtx.Logger.Level(zerzerolog.DebugLevel)
// 	ctx := context.Background()
//
// 	workerAgent, _ := New(bCtx)
// 	err := workerAgent.Connect(bCtx)
// 	assert.NoError(t, err)
// 	tests := []struct {
// 		name    string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Test worker component Run",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := workerAgent.Run(bCtx, ctx); (err != nil) != tt.wantErr {
// 				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
//
// }
