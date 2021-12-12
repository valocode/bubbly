package integrations

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/artifact/local"
	"github.com/aquasecurity/fanal/cache"
	"github.com/stretchr/testify/require"
)

func TestFanal(t *testing.T) {
	// c, err := cache.NewFSCache(os.TempDir())
	// require.NoError(t, err)

	// appl := applier.NewApplier(c)
	// appl.ApplyLayers()

	// anal := analyzer.NewAnalyzer([]analyzer.Type{})
	// var wg sync.WaitGroup
	// var result analyzer.AnalysisResult
	// limit := semaphore.NewWeighted(3)
	// anal.AnalyzeFile(context.TODO(), wg, limit, &result, "../ui", "", )
	fsCache, err := cache.NewFSCache(os.TempDir())
	require.NoError(t, err)
	a, err := local.NewArtifact("../docs", fsCache, artifact.Option{}, config.ScannerOption{})
	require.NoError(t, err)
	ref, err := a.Inspect(context.TODO())
	require.NoError(t, err)

	t.Log("name: ", ref.Name)
	t.Log("type: ", ref.Type)
	t.Log("blobs: ", ref.BlobIDs)
	// spew.Dump(ref.ImageMetadata)

	blobInfo, err := fsCache.GetBlob(ref.BlobIDs[0])
	require.NoError(t, err)
	for _, app := range blobInfo.Applications {
		// if app.FilePath != "package.json" {
		// 	continue
		// }
		fmt.Println("app: ", app.Type, app.FilePath)
		fmt.Println("libs: ", len(app.Libraries))
		// for _, lib := range app.Libraries {
		// 	fmt.Println("\tlib: ", lib.Name, lib.Version)
		// }
	}
	// blobInfo.Applications[0].Libraries[0].
	// fmt.Println("infos: ", len(blobInfo.PackageInfos))
	// for _, info := range blobInfo.PackageInfos {
	// 	fmt.Println("pkgs: ", len(info.Packages))
	// 	for _, pkg := range info.Packages {
	// 		fmt.Println("pkg: ", pkg.Name, pkg.License)
	// 	}
	// }
	// spew.Dump(blobInfo)

	// fapplier := applier.NewApplier(fsCache)
	// details, err := fapplier.ApplyLayers(ref.ImageMetadata.ID, ref.BlobIDs)
	// require.NoError(t, err)
	// spew.Dump(details)

}
