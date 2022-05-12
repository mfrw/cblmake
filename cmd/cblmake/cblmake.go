package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/mfrw/cblmake/internal/exe"
	"github.com/mfrw/cblmake/pkg/rpmexpander"

	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/graph/grapher"
	gpkgfetcher "github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/graph/pkgfetcher"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/graph/pkggraph"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/image/imager"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/image/isomaker"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/logger"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/scheduler"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/specreader"
	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/srpmpacker"
)

const (
	SOURCE_URL      = "https://cblmarinerstorage.blob.core.windows.net/sources/core"
	defaultBuildDir = "./build/SRPMS"
)

var (
	app       = kingpin.New("cblmake", "A tool to generate iso|vhd[x] from a custom input spec(s) dir")
	logFile   = exe.LogFileFlag(app)
	logLevel  = exe.LogLevelFlag(app)
	specsDir  = app.Flag("specs", "Path to the SPEC directory to create SRPMs from.").Required().String()
	distTag   = app.Flag("dist-tag", "The distribution tag SRPMs will be built with.").Default(".cm2").String()
	sourceURL = app.Flag("source-url", "URL to a source server to download SPEC sources from.").Default(SOURCE_URL).String()
	buildDir  = app.Flag("build-dir", "Directory to store temporary files while building.").Default(defaultBuildDir).String()
)

var (
	// TODO: embed the worker-tar using 'embed' from golang std lib
	// Something like this:
	//
	//      import "embed"
	//
	//      //go:embed chroot.tar.gz
	//      var f embed.FS
	//      chroot, _ := f.ReadFile("chroot.tar.gz")
	workerTar = app.Flag("worker-tar", "Full path to worker_chroot.tar.gz.  If this argument is empty, specs will be parsed in the host environment.").ExistingFile()
)

func main() {
	// TODO: should we cobra instead of kingpin? kingpin's last commit in 2021
	// cobra seems to be more feature complete and with more stars ;)
	// star count:
	// kingpin -> 3.2k
	// cobra   -> 26.6k

	// TODO: Think more on the logger.
	// Although `logrus` is cool to build upon, but should we maybe look at `uber/zap` ?

	app.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))
	logger.InitBestEffort(*logFile, *logLevel)

	// TODO: Remove hardcoded values. This is gross :(
	intSpecs := "/home/mfrw/cblmakedemo/tmp/cblmake/out/specs"
	intSrpms := "/home/mfrw/cblmakedemo/tmp/cblmake/out/srpms"
	specRdrOutJson := "/home/mfrw/cblmakedemo/tmp/cblmake/out/specreader.json"
	grapherOutput := "/home/mfrw/cblmakedemo/tmp/cblmake/out/grapher.dot"
	resolvedGraph := "/home/mfrw/cblmakedemo/tmp/cblmake/out/resolved.dot"

	// Step 1: Pack all the spec-files into SRPMs from the custom-specs dir
	packSrpms()

	// Step 2: Extract the SRPMs created above into specs
	extractSRPMS(intSpecs, intSrpms, 12)

	// Step 3: Read All the SPECS and create a flat json
	specReader(intSpecs, intSrpms, specRdrOutJson)

	// Step 4: Slurp the json file and create dependency graph
	createGraph(specRdrOutJson, grapherOutput)

	// Step 5: Slurp the dependency graph & fetch/resolve local/upstream dependencies
	resolvePackages(grapherOutput, resolvedGraph)

	// TODO (mfrw): Complete the below functions to have a full functional demo

	// Step 6: Build all the local packages
	buildPackages()

	// Step 7: Create a VHD/ISO
	buildImage()
	buildISO()

	// Step 8:
	// Cleanup
}

func buildPackages() {
	//TODO: Fill cfg with proper values
	// cfg := &scheduler.Config{}
	cfg := (*scheduler.Config)(nil)
	if cfg == nil {
		logger.Log.Info("TODO: Implement")
		return
	}
	scheduler.ScheduleBuild(cfg)
}

func buildImage() {
	//TODO: Fill cfg with proper values
	//cfg := &imager.Config{}
	cfg := (*imager.Config)(nil)
	if cfg == nil {
		logger.Log.Info("TODO: Implement")
		return
	}
	imager.CreateImage(cfg)
}

func buildISO() {
	//TODO: Fill isoMaker  with proper values
	isoMaker := (*isomaker.IsoMaker)(nil)
	if isoMaker == nil {
		logger.Log.Info("TODO: Implement")
		return
	}
	isoMaker.Make()
}

func resolvePackages(in, out string) {
	baseDir := "/home/mfrw/cblmakedemo/tmp/cblmake"
	cfg := &gpkgfetcher.Config{
		InputGraph:        in,
		OutputGraph:       out,
		OutDir:            "/home/mfrw/cblmakedemo/tmp/cblmake/out/downloaded_rpms",
		ExistingRpmDir:    "/home/mfrw/cblmakedemo/tmp/cblmake/rpms",
		TmpDir:            "/home/mfrw/cblmakedemo/tmp/cblmake/tmp",
		WorkerTar:         filepath.Join(baseDir, "chroot.tar.gz"),
		RepoFiles:         []string{filepath.Join(baseDir, "cm.repo")},
		OutputSummaryFile: filepath.Join(baseDir, "download_summary"),
	}
	gpkgfetcher.ResolvePackages(cfg)
}

func createGraph(input, output string) *pkggraph.PkgGraph {
	cfg := &grapher.Config{
		Input:  input,
		Output: output,
	}
	depGraph, _ := grapher.GenerateDependencyGraph(cfg)
	pkggraph.WriteDOTGraphFile(depGraph, output)
	return depGraph
}

func packSrpms() error {
	baseDir := "/home/mfrw/cblmakedemo/tmp/cblmake"
	cfg := &srpmpacker.Config{
		SpecsDir:          *specsDir,
		OutDir:            filepath.Join(baseDir, "out", "srpms"),
		BuildDir:          filepath.Join(baseDir, "build-dir"),
		DistTag:           *distTag,
		WorkerTar:         filepath.Join(baseDir, "chroot.tar.gz"),
		SignatureHandling: "update",
		SourceURL:         *sourceURL,
		Workers:           12,
	}
	return srpmpacker.CreateAllSRPMsWrapper(cfg)
}

func extractSRPMS(dstDir, srcDir string, nrWorkers int) {
	semaphore := make(chan struct{}, nrWorkers)
	files, _ := ioutil.ReadDir(srcDir)
	for _, v := range files {
		if strings.HasSuffix(v.Name(), ".src.rpm") {
			srcRpm := filepath.Join(srcDir, v.Name())
			dstDir := filepath.Join(dstDir, strings.TrimSuffix(v.Name(), ".src.rpm"))
			semaphore <- struct{}{}
			go func(dstDir, srcRpm string) {
				rpmexpander.ExtractRPM(dstDir, srcRpm)
				<-semaphore
			}(dstDir, srcRpm)
		}
	}
	close(semaphore)
}

func specReader(specsDir, srpmsDir, outDir string) {
	baseDir := "/home/mfrw/cblmakedemo/tmp/cblmake"
	cfg := &specreader.Config{
		SpecsDir:  specsDir,
		SrpmsDir:  srpmsDir,
		Output:    outDir,
		BuildDir:  filepath.Join(baseDir, "build-dir", "specreader"),
		RpmsDir:   filepath.Join(baseDir, "out", "rpms"),
		DistTag:   *distTag,
		WorkerTar: filepath.Join(baseDir, "chroot.tar.gz"),
		Workers:   12,
	}
	specreader.ParseSPECsWrapper(cfg)
}
