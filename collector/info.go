package collector

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/template"

	"github.com/prometheus/client_golang/prometheus"
)

// Build information. Populated at build-time.
var (
	BuildVersion string
	BuildBranch  string
	BuildUser    string
	BuildDate    string
	GoVersion    = runtime.Version()
	GoOS         = runtime.GOOS
	GoArch       = runtime.GOARCH
)

// NewCollector returns a collector that exports metrics about current version
// information.
func VersionCollector(program, BuildVersion, BuildBranch string) prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: program,
			Name:      "build_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by version, revision, branch, goversion from which %s was built, and the goos and goarch for the build.",
				program,
			),
			ConstLabels: prometheus.Labels{
				"version":   BuildVersion,
				"branch":    BuildBranch,
				"goversion": GoVersion,
				"goos":      GoOS,
				"goarch":    GoArch,
			},
		},
		func() float64 { return 1 },
	)
}

// versionInfoTmpl contains the template used by Info.
var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
  platform:         {{.platform}}
  tags:             {{.tags}}
`

// Print returns version information.
func VersionPrint(program, BuildVersion, BuildUser, BuildDate, BuildBranch string) string {
	m := map[string]string{
		"program":   program,
		"version":   BuildVersion,
		"branch":    BuildBranch,
		"buildUser": BuildUser,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
		"platform":  GoOS + "/" + GoArch,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
