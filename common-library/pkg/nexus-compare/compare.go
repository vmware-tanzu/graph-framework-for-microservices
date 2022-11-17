package nexus_compare

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	SpecMatch   = ".*properties\\/spec.*"
	StatusMatch = ".*properties\\/status.*"
)

func CompareDirs(dir1, dir2 string) (bool, *bytes.Buffer) {

	ans := false
	buffer := new(bytes.Buffer)
	files1, err := os.ReadDir(dir1)
	if err != nil {
		log.Fatal(err)
	}
	files2, err := os.ReadDir(dir2)
	if err != nil {
		log.Fatal(err)
	}

	for _, f1 := range files1 {
		for _, f2 := range files2 {
			if f1.Name() == f2.Name() && f1.IsDir() {
				a, b := CompareDirs(path.Join(dir1, f1.Name()), path.Join(dir2, f2.Name()))
				a = ans || a
				buffer.Write(b.Bytes())
			} else if f1.Name() == f2.Name() && strings.HasSuffix(f1.Name(), ".yaml") {
				p1 := path.Join(dir1, f1.Name())
				p2 := path.Join(dir2, f2.Name())
				a, b := CompareFiles(p1, p2)
				ans = ans || a
				buffer.Write(b.Bytes())
			}
		}
	}

	if ans == false {
		buffer.WriteString("no changes detected")
	}

	return ans, buffer
}

func CompareFiles(file1, file2 string) (bool, *bytes.Buffer) {
	buffer := new(bytes.Buffer)
	headerColor, _ := colorful.Hex("#B9311B")
	fileColor, _ := colorful.Hex("#088F8F")

	spec, status, nexus := GetCompareReports(file1, file2)
	spd := len(spec.Diffs) > 0
	std := len(status.Diffs) > 0
	nd := len(nexus.Diffs) > 0

	if !spd && !std && !nd {
		return false, buffer
	}

	_, _ = buffer.WriteString(bunt.Style(
		"detected changes in model stored in ",
		bunt.EachLine(),
		bunt.Foreground(headerColor),
	))
	_, _ = buffer.WriteString(bunt.Style(
		file1,
		bunt.EachLine(),
		bunt.Foreground(fileColor),
	))
	_, _ = buffer.WriteString("\n\n")

	if spd {
		buffer.WriteString("spec changes: ")
		PrintReportDiff(spec, buffer)
	}
	if spd {
		buffer.WriteString("status changes: ")
		PrintReportDiff(status, buffer)
	}
	if spd {
		buffer.WriteString("nexus annotation changes: ")
		PrintReportDiff(nexus, buffer)
	}

	return true, buffer
}

func GetCompareReports(file1, file2 string) (dyff.Report, dyff.Report, dyff.Report) {
	from, to, err := ytbx.LoadFiles(file1, file2)
	if err != nil {
		fmt.Println(err, "failed to load input files")
	}

	report, err := dyff.CompareInputFiles(from, to,
		dyff.IgnoreOrderChanges(true),
		dyff.KubernetesEntityDetection(true),
		dyff.AdditionalIdentifiers(""),
	)

	if err != nil {
		fmt.Println(err)
	}

	specDiffs := getSpecificReport(report, SpecMatch)
	statusDiffs := getSpecificReport(report, StatusMatch)
	nexusDiffs := getAnnotationReport(report)

	return specDiffs, statusDiffs, nexusDiffs

}

func getSpecificReport(r dyff.Report, match string) dyff.Report {
	var res []dyff.Diff
	for _, d := range r.Diffs {

		cc, _ := regexp.Match(match, []byte(d.Path.String()))
		if cc {
			res = append(res, d)
		}
	}
	r.Diffs = res
	return r
}

func getNexusAnnotation(file ytbx.InputFile) string {
	return file.Documents[0].Content[0].Content[5].Content[1].Content[1].Value
}

func getAnnotationReport(r dyff.Report) dyff.Report {
	aNexus := getNexusAnnotation(r.From)
	bNexus := getNexusAnnotation(r.To)

	aFile, err := os.CreateTemp("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	bFile, err := os.CreateTemp("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(aFile.Name())
	defer os.Remove(bFile.Name())
	aFile.Write([]byte(aNexus))
	bFile.Write([]byte(bNexus))

	from, to, err := ytbx.LoadFiles(aFile.Name(), bFile.Name())
	if err != nil {
		fmt.Println(err, "failed to load input files")
	}

	report, err := dyff.CompareInputFiles(from, to,
		dyff.IgnoreOrderChanges(true),
		dyff.KubernetesEntityDetection(true),
		dyff.AdditionalIdentifiers(""),
	)

	return report

}

func PrintReportDiff(report dyff.Report, buffer *bytes.Buffer) {
	h := &dyff.HumanReport{
		Report:               report,
		DoNotInspectCerts:    true,
		NoTableStyle:         true,
		OmitHeader:           true,
		UseGoPatchPaths:      true,
		MinorChangeThreshold: 0.1,
	}
	if err := h.WriteReport(buffer); err != nil {
		fmt.Println(err)
		return
	}
}
