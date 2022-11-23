package nexus_compare

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

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

	spec, status, nexus, err := CompareReports(data1, data2)
	if err != nil {
		return true, nil, err
	}

	spd := len(spec.Diffs) > 0
	std := len(status.Diffs) > 0
	nd := len(nexus.Diffs) > 0
	if !spd && !std && !nd {
		return false, buffer, nil
	}

	_, err = buffer.WriteString(bunt.Style(
		"detected changes in model stored in ",
		bunt.EachLine(),
		bunt.Foreground(headerColor),
	))
	if err != nil {
		return true, nil, err
	}
	_, err = buffer.WriteString(bunt.Style(
		name,
		bunt.EachLine(),
		bunt.Foreground(fileColor),
	))
	if err != nil {
		return true, nil, err
	}
	_, err = buffer.WriteString("\n\n")
	if err != nil {
		return true, nil, err
	}

	if spd {
		_, err = buffer.WriteString("spec changes: ")
		if err != nil {
			return true, nil, err
		}
		err := PrintReportDiff(spec, buffer)
		if err != nil {
			return true, nil, err
		}
	}
	if err != nil {
		return true, nil, err
	}

	if std {
		_, err = buffer.WriteString("status changes: ")
		if err != nil {
			return true, nil, err
		}
		err = PrintReportDiff(status, buffer)
		if err != nil {
			return true, nil, err
		}
	}
	if nd {
		_, err = buffer.WriteString("nexus annotation changes: ")
		if err != nil {
			return true, nil, err
		}
		err = PrintReportDiff(nexus, buffer)
		if err != nil {
			return true, nil, err
		}
	}

	return true, buffer, nil
}

func CompareReports(data1, data2 []byte) (dyff.Report, dyff.Report, dyff.Report, error) {
	file1, err := writeTempFile(data1)
	if err != nil {
		return dyff.Report{}, dyff.Report{}, dyff.Report{}, err
	}
	file2, err := writeTempFile(data2)
	if err != nil {
		return dyff.Report{}, dyff.Report{}, dyff.Report{}, err
	}
	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	from, to, err := ytbx.LoadFiles(file1.Name(), file2.Name())
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

	nexusDiffs, err := getAnnotationReport(data1, data2)
	if err != nil {
		return dyff.Report{}, dyff.Report{}, dyff.Report{}, err
	}

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
	var diffs []dyff.Diff
	for _, di := range r.Diffs {
		var ds []dyff.Detail
		for _, d := range di.Details {
			if d.From != nil || d.To == nil {
				ds = append(ds, d)
			}
		}
		if ds == nil {
			continue
		}
		di.Details = ds
		diffs = append(diffs, di)
	}
	r.Diffs = diffs
	return r
}

func getAnnotationReport(data1, data2 []byte) (dyff.Report, error) {
	aNexus, err := getMapNode(data1, []string{"metadata", "annotations", "nexus"})
	if err != nil {
		return dyff.Report{}, err
	}
	bNexus, err := getMapNode(data2, []string{"metadata", "annotations", "nexus"})
	if err != nil {
		return dyff.Report{}, err
	}

	aFile, err := writeTempFile([]byte(aNexus.(string)))
	if err != nil {
		return dyff.Report{}, err
	}
	bFile, err := writeTempFile([]byte(bNexus.(string)))
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
	h := CustomReport{
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

func getMapNode(data []byte, path []string) (interface{}, error) {
	var t interface{}
	var ok bool
	errPath := "root"
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return map[string]interface{}{}, err
	}
	for _, p := range path {
		errPath += fmt.Sprintf(".%s", p)
		t, ok = t.(map[string]interface{})[p]
		if !ok {
			return map[string]interface{}{}, errors.New(fmt.Sprintf("%s not found while looking for nexus annotation", errPath))
		}
	}
	return t, err
}
