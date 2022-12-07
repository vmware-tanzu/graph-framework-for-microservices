package dir

import (
	"errors"
	"fmt"

	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/kube-wrapper"
)

func ApplyDir(directory string, force bool, c kubewrapper.ClientInt, cFunc compareFunc) error {
	// check for incompatible models and not installed. Return  if any and force != true
	incNames, _, text, err := CheckDir(directory, c, cFunc)
	if err != nil {
		return err
	}
	if len(incNames) > 0 && force == false {
		fmt.Println(text)
		return errors.New("changes detected. If you want to install models anyway, run with -force=true")
	}

	// check if any data for incompatible models and return if so
	var dataExist []string
	for _, n := range incNames {
		res, err := c.ListResources(*c.GetCrd(n))
		fmt.Println(res)
		if err != nil {
			return err
		}
		if len(res) > 0 {
			dataExist = append(dataExist, n)
		}
	}
	if len(dataExist) > 0 {
		return errors.New(fmt.Sprintf("There are some data that exist in datamodels that are backward incompatible: %v. Please remove them manually to force upgrade CRDs", dataExist))
	}

	// upsert all the models
	err = InstallDir(directory, c)
	if err != nil {
		return err
	}
	return nil
}
