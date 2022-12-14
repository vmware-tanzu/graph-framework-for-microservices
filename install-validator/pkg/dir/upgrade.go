package dir

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
)

// ApplyDir checks if there are any incompatible crds and data for them. Applies them based on force flag.
func ApplyDir(directory string, force bool, c kubewrapper.ClientInt, cFunc compareFunc) error {
	err := c.FetchCrds()
	if err != nil {
		return err
	}

	// check for incompatible and outdated models. Return if there are any incompatibilities and force != true
	inCompatibleCRDs, outdated, err := CheckDir(directory, c, cFunc)
	if err != nil {
		return err
	}
	if (len(inCompatibleCRDs) > 0 || len(outdated) > 0) && !force {
		textChanges := new(bytes.Buffer)
		for _, txt := range inCompatibleCRDs {
			textChanges.Write(txt.Bytes())
		}
		for _, txt := range outdated {
			textChanges.WriteString(fmt.Sprintf("%s is deleted\n", txt))
		}
		logrus.Warn(textChanges)
		return errors.New("incompatible datamodel changes detected")
	}

	// check if there are any resources for incompatible crds and return if so
	var incompatibleAll []string //changed and deleted crds should both be checked
	for crd := range inCompatibleCRDs {
		incompatibleAll = append(incompatibleAll, crd)
	}
	incompatibleAll = append(incompatibleAll, outdated...)
	var cr []string
	for _, crd := range incompatibleAll {
		res, err := c.ListResources(*c.GetCrd(crd))
		if err != nil {
			return err
		}
		if len(res) > 0 {
			cr = append(cr, crd)
		}
	}
	if len(cr) > 0 {
		return fmt.Errorf("validation failed as objects exists in the system for the incompatibly changed nodes: %v", cr)
	}

	//delete outdated models
	for _, o := range outdated {
		err = c.DeleteCrd(o)
		if err != nil {
			return err
		}
	}

	// upsert all the models
	err = InstallDir(directory, c)
	if err != nil {
		return err
	}
	return nil
}
