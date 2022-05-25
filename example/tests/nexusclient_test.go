package crd_generator_test

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	configv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/config.tsm.tanzu.vmware.com/v1"
	gnsv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/gns.tsm.tanzu.vmware.com/v1"
	rootv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/root.tsm.tanzu.vmware.com/v1"
	sgv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"
	nexus_client "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/nexus-client"
)

var _ = Describe("Template renderers tests", func() {
	var (
		//err error
		fakeClient *nexus_client.Clientset
		str        gnsv1.MyStr = "test"
	)
	BeforeEach(func() {
		fakeClient = nexus_client.NewFakeClient()
	})

	It("should create root object", func() {
		rootDef := &rootv1.Root{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		}
		root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
		Expect(err).NotTo(HaveOccurred())
		Expect(root.DisplayName()).To(Equal("default"))
		// expect name to be hashed
		Expect(root.GetName()).To(Equal("9d336ed798cf54e3ef224fb00017b75b1a15abff"))
	})

	Context("Child objects", func() {
		var (
			err  error
			root *nexus_client.RootRoot
		)

		BeforeEach(func() {
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			root, err = fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create child object", func() {
			cfgName := "configObj"
			expectedLabels := map[string]string{
				"roots.root.tsm.tanzu.vmware.com": "default",
				"nexus/display_name":              "configObj",
				"nexus/is_name_hashed":            "true",
			}
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal("configObj"))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr).To(Equal(&str))

			// GetConfig should return same object
			cfg, err = root.GetConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr).To(Equal(&str))

			// Also Get by using hashed name should return same thing
			cfg, err = fakeClient.Config().GetConfigByName(context.TODO(), cfg.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr).To(Equal(&str))

			// Another create should fail
			_, err = root.AddConfig(context.TODO(), cfgDef)
			Expect(err).To(HaveOccurred())
		})

		It("should delete child", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))

			err = root.DeleteConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			cfg, err = root.GetConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).To(BeNil())
		})

		It("should update spec of object", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(*cfg.Spec.MyStr).To(Equal(str))

			var updatedStr gnsv1.MyStr = "updatedStr"
			cfg.Spec.MyStr = &updatedStr
			err = cfg.Update(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Spec.MyStr).To(Equal(&updatedStr))
			cfg, err = root.GetConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Spec.MyStr).To(Equal(&updatedStr))
		})

		It("should create named child", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			gnsDef := &gnsv1.Gns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gnsName",
				},
			}
			gns, err := cfg.AddGNS(context.TODO(), gnsDef)
			Expect(err).NotTo(HaveOccurred())
			sg1def := &sgv1.SvcGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "sg1",
				},
			}
			sg1, err := gns.AddGnsServiceGroups(context.TODO(), sg1def)
			Expect(err).NotTo(HaveOccurred())
			Expect(sg1.DisplayName()).To(Equal("sg1"))
			sg2def := &sgv1.SvcGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "sg2",
				},
			}
			sg2, err := gns.AddGnsServiceGroups(context.TODO(), sg2def)
			Expect(err).NotTo(HaveOccurred())
			Expect(sg2.DisplayName()).To(Equal("sg2"))

			getSg1, err := gns.GetGnsServiceGroups(context.TODO(), "sg1")
			Expect(err).NotTo(HaveOccurred())
			Expect(getSg1.DisplayName()).To(Equal("sg1"))

			getSg2, err := gns.GetGnsServiceGroups(context.TODO(), "sg2")
			Expect(err).NotTo(HaveOccurred())
			Expect(getSg2.DisplayName()).To(Equal("sg2"))

			allSgs, err := gns.GetAllGnsServiceGroups(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(allSgs).To(HaveLen(2))

			listSgs, err := fakeClient.Servicegroup().ListSvcGroups(context.TODO(), metav1.ListOptions{})
			Expect(listSgs).To(HaveLen(2))
		})

		It("should remove all children when parent is removed", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			cfg, err = fakeClient.Config().GetConfigByName(context.TODO(), cfg.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))
			gnsDef := &gnsv1.Gns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gnsName",
				},
			}
			gns, err := cfg.AddGNS(context.TODO(), gnsDef)
			Expect(err).NotTo(HaveOccurred())
			gns, err = fakeClient.Gns().GetGnsByName(context.TODO(), gns.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(gns.DisplayName()).To(Equal("gnsName"))

			err = root.Delete(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			cfg, err = fakeClient.Config().GetConfigByName(context.TODO(), cfg.GetName())
			Expect(err).To(HaveOccurred())
			Expect(cfg).To(BeNil())
			Expect(err).To(HaveOccurred())

			gns, err = fakeClient.Gns().GetGnsByName(context.TODO(), gns.GetName())
			Expect(err).To(HaveOccurred())
			Expect(gns).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Linking objects", func() {
		fakeClient = nexus_client.NewFakeClient()
		rootDef := &rootv1.Root{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		}
		root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
		Expect(err).NotTo(HaveOccurred())
		cfgDef := &configv1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cfg",
			},
			Spec: configv1.ConfigSpec{
				MyStr: &str,
			},
		}
		cfg, err := root.AddConfig(context.TODO(), cfgDef)
		Expect(err).NotTo(HaveOccurred())
		gnsDef := &gnsv1.Gns{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gnsName",
			},
		}
		gns, err := cfg.AddGNS(context.TODO(), gnsDef)
		Expect(err).NotTo(HaveOccurred())
		dnsDef := &gnsv1.Dns{
			ObjectMeta: metav1.ObjectMeta{
				Name: "dnsName",
			},
		}
		dns, err := gns.GetDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(dns).To(BeNil())

		dns, err = cfg.AddDNS(context.TODO(), dnsDef)
		Expect(err).NotTo(HaveOccurred())

		err = gns.LinkDns(context.TODO(), dns)
		Expect(err).NotTo(HaveOccurred())

		getLinkedDns, err := gns.GetDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(getLinkedDns.DisplayName()).To(Equal("dnsName"))

		err = gns.UnlinkDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())

		getLinkedDns, err = gns.GetDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(getLinkedDns).To(BeNil())
	})

	Context("Getting parent", func() {
		fakeClient = nexus_client.NewFakeClient()
		rootDef := &rootv1.Root{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		}
		root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
		Expect(err).NotTo(HaveOccurred())
		cfgDef := &configv1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cfg",
			},
			Spec: configv1.ConfigSpec{
				MyStr: "test",
			},
		}
		cfg, err := root.AddConfig(context.TODO(), cfgDef)
		Expect(err).NotTo(HaveOccurred())
		gnsDef := &gnsv1.Gns{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gnsName",
			},
		}
		gns, err := cfg.AddGNS(context.TODO(), gnsDef)
		Expect(err).NotTo(HaveOccurred())

		gnsParent, err := gns.GetParent(context.TODO())
		Expect(err).NotTo(HaveOccurred())

		Expect(gnsParent.DisplayName()).To(Equal("cfg"))

		configParent, err := gnsParent.GetParent(context.TODO())
		Expect(configParent.DisplayName()).To(Equal("default"))
	})
})
