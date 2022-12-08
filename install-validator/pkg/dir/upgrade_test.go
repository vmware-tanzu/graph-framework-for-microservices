package dir

import (
	"bytes"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
	mockkubernetes "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils/mocks"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestApplyDir(t *testing.T) {
	type args struct {
		directory string
		force     bool
		c         kubewrapper.ClientInt
		cFunc     compareFunc
	}
	ctrl := gomock.NewController(t)
	cc := mockkubernetes.NewMockClientInt(ctrl)

	patt1f, err := os.ReadFile("./test_dir/patt1.yaml")
	if err != nil {
		t.Fatal("error while init", err)
	}
	rootf, _ := os.ReadFile("./test_dir2/root_root.yaml")
	if err != nil {
		t.Fatal("error while init", err)
	}
	rootOutdatedf, _ := os.ReadFile("./test_dir/root_root_outdated.yaml")
	if err != nil {
		t.Fatal("error while init", err)
	}
	var patt1 v1.CustomResourceDefinition
	var rootRoot v1.CustomResourceDefinition
	var rootOutdated v1.CustomResourceDefinition
	err = yaml.Unmarshal(patt1f, &patt1)
	if err != nil {
		t.Fatal("error while init", err)
	}
	err = yaml.Unmarshal(rootf, &rootRoot)
	if err != nil {
		t.Fatal("error while init", err)
	}
	err = yaml.Unmarshal(rootOutdatedf, &rootOutdated)
	if err != nil {
		t.Fatal("error while init", err)
	}

	cc.EXPECT().FetchCrds().Return(nil).AnyTimes()

	cc.EXPECT().ApplyCrd(patt1).Return(nil).AnyTimes()
	cc.EXPECT().ApplyCrd(rootRoot).Return(nil).AnyTimes()
	cc.EXPECT().ApplyCrd(rootOutdated).Return(nil).AnyTimes()

	cc.EXPECT().GetCrd("my-crds.com.example").Return(&patt1).AnyTimes()
	cc.EXPECT().GetCrd("roots.rootoutdated.tsm.tanzu.vmware.com").Return(&patt1).AnyTimes()
	cc.EXPECT().GetCrd("roots.root.tsm.tanzu.vmware.com").Return(&rootRoot).AnyTimes()

	cc.EXPECT().ListResources(patt1).Return([]interface{}{"aa"}, nil).AnyTimes()
	cc.EXPECT().ListResources(rootRoot).Return([]interface{}{}, nil).AnyTimes()

	cc.EXPECT().GetCrds().Return([]v1.CustomResourceDefinition{}).AnyTimes()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all new data",
			args: args{
				directory: "./test_dir",
				force:     false,
				c:         cc,
				cFunc:     getNoDiff,
			},
			wantErr: false,
		},
		{
			name: "diffs no force",
			args: args{
				directory: "./test_dir",
				force:     false,
				c:         cc,
				cFunc:     AnyDiffs,
			},
			wantErr: true,
		},
		{
			name: "diffs force, data exist",
			args: args{
				directory: "./test_dir",
				force:     true,
				c:         cc,
				cFunc:     AnyDiffs,
			},
			wantErr: true,
		},
		{
			name: "diffs force, no data",
			args: args{
				directory: "./test_dir2",
				force:     true,
				c:         cc,
				cFunc:     AnyDiffs,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ApplyDir(tt.args.directory, tt.args.force, tt.args.c, tt.args.cFunc); (err != nil) != tt.wantErr {
				t.Errorf("ApplyDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getNoDiff(_ []byte, _ []byte) (bool, *bytes.Buffer, error) {
	return false, nil, nil
}
func AnyDiffs(_ []byte, _ []byte) (bool, *bytes.Buffer, error) {
	return true, bytes.NewBuffer(nil), nil
}
