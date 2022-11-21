package nexus_compare

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	SpecMatch   = ".*properties\\/spec.*"
	StatusMatch = ".*properties\\/status.*"
)

func CompareFiles(data1, data2 []byte) (bool, *bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	headerColor, _ := colorful.Hex("#B9311B")
	fileColor, _ := colorful.Hex("#088F8F")

	name, err := getSpecName(data1)
	if err != nil {
		return true, nil, err
	}

	file1, err := writeTempFile(data1)
	if err != nil {
		return true, nil, err
	}
	file2, err := writeTempFile(data2)
	if err != nil {
		return true, nil, err
	}
	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	spec, status, nexus, err := CompareReports(file1.Name(), file2.Name())
	if err != nil {
		return true, nil, err
	}

	spd := len(spec.Diffs) > 0
	std := len(status.Diffs) > 0
	nd := len(nexus.Diffs) > 0
	if !spd && !std && !nd {
		return true, buffer, nil
	}

	_, _ = buffer.WriteString(bunt.Style(
		"detected changes in model stored in ",
		bunt.EachLine(),
		bunt.Foreground(headerColor),
	))
	_, _ = buffer.WriteString(bunt.Style(
		name,
		bunt.EachLine(),
		bunt.Foreground(fileColor),
	))
	_, _ = buffer.WriteString("\n\n")
	if spd {
		buffer.WriteString("spec changes: ")
		err := PrintReportDiff(spec, buffer)
		if err != nil {
			return true, nil, err
		}
	}
	if spd {
		buffer.WriteString("status changes: ")
		err := PrintReportDiff(status, buffer)
		if err != nil {
			return true, nil, err
		}
	}
	if spd {
		buffer.WriteString("nexus annotation changes: ")
		err := PrintReportDiff(nexus, buffer)
		if err != nil {
			return true, nil, err
		}
	}

	return true, buffer, nil
}

func CompareReports(file1, file2 string) (dyff.Report, dyff.Report, dyff.Report, error) {
	from, to, err := ytbx.LoadFiles(file1, file2)
	if err != nil {
		return dyff.Report{}, dyff.Report{}, dyff.Report{}, err
	}

	report, err := dyff.CompareInputFiles(from, to,
		dyff.IgnoreOrderChanges(true),
		dyff.KubernetesEntityDetection(true),
		dyff.AdditionalIdentifiers(""),
	)

	if err != nil {
		return dyff.Report{}, dyff.Report{}, dyff.Report{}, err
	}

	sr := getSpecificReport(report, SpecMatch)
	sd := getSpecificReport(report, StatusMatch)
	specDiffs := filterReport(&sr)
	statusDiffs := filterReport(&sd)
	nexusDiffs, err := getAnnotationReport(report)

	return *specDiffs, *statusDiffs, nexusDiffs, err

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

func filterReport(r *dyff.Report) *dyff.Report {
	for i, di := range r.Diffs {
		var ds []dyff.Detail
		for _, d := range di.Details {
			if d.From != nil || d.To == nil {
				ds = append(ds, d)
			}
		}
		di.Details = ds
		r.Diffs[i] = di
	}
	return r
}

func getNexusAnnotation(file ytbx.InputFile) string {
	return file.Documents[0].Content[0].Content[5].Content[1].Content[1].Value
}

func getAnnotationReport(r dyff.Report) (dyff.Report, error) {
	aNexus := getNexusAnnotation(r.From)
	bNexus := getNexusAnnotation(r.To)

	aFile, err := writeTempFile([]byte(aNexus))
	if err != nil {
		return dyff.Report{}, err
	}
	bFile, err := writeTempFile([]byte(bNexus))
	if err != nil {
		return dyff.Report{}, err
	}
	defer os.Remove(aFile.Name())
	defer os.Remove(bFile.Name())

	from, to, err := ytbx.LoadFiles(aFile.Name(), bFile.Name())
	if err != nil {
		return dyff.Report{}, err
	}

	report, err := dyff.CompareInputFiles(from, to,
		dyff.IgnoreOrderChanges(true),
		dyff.KubernetesEntityDetection(true),
		dyff.AdditionalIdentifiers(""),
	)

	return report, err

}

func PrintReportDiff(report dyff.Report, buffer *bytes.Buffer) error {
	h := &dyff.HumanReport{
		Report:               report,
		DoNotInspectCerts:    true,
		NoTableStyle:         true,
		OmitHeader:           true,
		UseGoPatchPaths:      true,
		MinorChangeThreshold: 0.1,
	}
	if err := h.WriteReport(buffer); err != nil {
		return err
	}
	return nil
}

func writeTempFile(data []byte) (*os.File, error) {
	aFile, err := os.CreateTemp("", "compare")
	_, err = aFile.Write(data)
	return aFile, err
}

func getSpecName(data []byte) (string, error) {
	t := make(map[string]interface{})
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return "", err
	}
	return t["metadata"].(map[string]interface{})["name"].(string), nil
}
