package nexus_compiler_test

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

var _ = Describe("Nexus clients tests", func() {
	var (
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
					MyStr0: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal("configObj"))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr0).To(Equal(&str))

			// GetConfig should return same object
			cfg, err = root.GetConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr0).To(Equal(&str))

			// Also Get by using hashed name should return same thing
			cfg, err = fakeClient.Config().GetConfigByName(context.TODO(), cfg.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))
			Expect(cfg.GetLabels()).To(BeEquivalentTo(expectedLabels))
			Expect(cfg.Spec.MyStr0).To(Equal(&str))

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
					MyStr0: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.DisplayName()).To(Equal(cfgName))

			gnsDef := &gnsv1.Gns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gnsName",
				},
			}
			gns, err := cfg.AddGNS(context.TODO(), gnsDef)
			Expect(err).NotTo(HaveOccurred())

			getGns, err := fakeClient.Gns().GetGnsByName(context.TODO(), gns.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(getGns.GetName()).To(Equal(gns.GetName()))

			err = root.DeleteConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			getGns, err = fakeClient.Gns().GetGnsByName(context.TODO(), gns.GetName())
			Expect(err).To(HaveOccurred())

			cfg, err = root.GetConfig(context.TODO())

			//Expect(err).NotTo(HaveOccurred())
			Expect(cfg).To(BeNil())
		})

		It("should update spec of object", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
				},
			}
			cfg, err := root.AddConfig(context.TODO(), cfgDef)
			Expect(err).NotTo(HaveOccurred())
			Expect(*cfg.Spec.MyStr0).To(Equal(str))

			var updatedStr gnsv1.MyStr = "updatedStr"
			cfg.Spec.MyStr0 = &updatedStr
			err = cfg.Update(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Spec.MyStr0).To(Equal(&updatedStr))
			cfg, err = root.GetConfig(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Spec.MyStr0).To(Equal(&updatedStr))
		})

		It("should create named child", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
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

			Expect(allSgs[0].DisplayName()).Should(Or(Equal("sg1"), Equal("sg2")))
			Expect(allSgs[1].DisplayName()).Should(Or(Equal("sg1"), Equal("sg2")))
			Expect(allSgs[0].DisplayName()).NotTo(Equal(allSgs[1].DisplayName()))

			listSgs, err := fakeClient.Servicegroup().ListSvcGroups(context.TODO(), metav1.ListOptions{})
			Expect(listSgs).To(HaveLen(2))

			Expect(listSgs[0].DisplayName()).Should(Or(Equal("sg1"), Equal("sg2")))
			Expect(listSgs[1].DisplayName()).Should(Or(Equal("sg1"), Equal("sg2")))
			Expect(listSgs[0].DisplayName()).NotTo(Equal(listSgs[1].DisplayName()))
		})

		It("should delete named child", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
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

			getSg2, err := fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), sg2.GetName())
			Expect(err).NotTo(HaveOccurred())
			Expect(getSg2.DisplayName()).To(Equal("sg2"))

			err = gns.DeleteGnsServiceGroups(context.TODO(), "sg2")
			Expect(err).NotTo(HaveOccurred())
			getSg1, err = gns.GetGnsServiceGroups(context.TODO(), "sg1")
			Expect(err).NotTo(HaveOccurred())
			Expect(getSg1.DisplayName()).To(Equal("sg1"))

			_, err = fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), getSg2.GetName())
			Expect(err).To(HaveOccurred())
		})

		It("should remove all children when parent is removed", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
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

			gns, err = fakeClient.Gns().GetGnsByName(context.TODO(), gns.GetName())
			Expect(err).To(HaveOccurred())
			Expect(gns).To(BeNil())
		})

		It("should delete all named children when parent is removed", func() {
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
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

			_, err = fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), sg1.GetName())
			Expect(err).NotTo(HaveOccurred())
			_, err = fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), sg2.GetName())
			Expect(err).NotTo(HaveOccurred())

			err = cfg.DeleteGNS(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			_, err = fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), sg1.GetName())
			Expect(err).To(HaveOccurred())
			_, err = fakeClient.Servicegroup().GetSvcGroupByName(context.TODO(), sg2.GetName())
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
				MyStr0: &str,
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
				Name: "default",
			},
		}
		dns, err := gns.GetDns(context.TODO())
		Expect(err).To(HaveOccurred())
		Expect(dns).To(BeNil())

		dns, err = cfg.AddDNS(context.TODO(), dnsDef)
		Expect(err).NotTo(HaveOccurred())

		err = gns.LinkDns(context.TODO(), dns)
		Expect(err).NotTo(HaveOccurred())

		getLinkedDns, err := gns.GetDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(getLinkedDns.DisplayName()).To(Equal("default"))

		err = gns.UnlinkDns(context.TODO())
		Expect(err).NotTo(HaveOccurred())

		getLinkedDns, err = gns.GetDns(context.TODO())
		Expect(err).To(HaveOccurred())
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
				MyStr0: &str,
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

	Context("Custom Errors", func() {
		It("should throw IsNotFound error when node's not present", func() {
			root, err := fakeClient.GetRootRoot(context.TODO())
			Expect(root).To(BeNil())
			Expect(nexus_client.IsNotFound(err)).To(BeTrue())
		})
		It("should throw ChildNotFound error when child is not present", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).NotTo(HaveOccurred())
			cfg, err := root.GetConfig(context.TODO())
			Expect(cfg).To(BeNil())
			Expect(err.Error()).To(Equal("child Config not found for Root.Root: default"))

			Expect(nexus_client.IsChildNotFound(err)).To(BeTrue())
			Expect(nexus_client.IsNotFound(err)).To(BeFalse())
		})

		It("should throw ChildNotFound error when named child is not present", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).NotTo(HaveOccurred())
			cfgName := "configObj"
			cfgDef := &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: cfgName,
				},
				Spec: configv1.ConfigSpec{
					MyStr0: &str,
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

			sgX, err := gns.GetGnsServiceGroups(context.TODO(), "sgX")
			Expect(sgX).To(BeNil())
			Expect(err.Error()).To(Equal("child GnsServiceGroups: sgX not found for Gns.Gns: gnsName"))

			Expect(nexus_client.IsChildNotFound(err)).To(BeTrue())
		})

		It("should throw LinkNotFound error when link is not present", func() {
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
					MyStr0: &str,
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

			dns, err := gns.GetDns(context.TODO())
			Expect(dns).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("link Dns not found for Gns.Gns: gnsName"))
			Expect(nexus_client.IsLinkNotFound(err)).To(BeTrue())
			Expect(nexus_client.IsNotFound(err)).To(BeFalse())

		})
	})

	Context("Singleton nodes", func() {
		It("should throw when user tries to create root singleton object which doesn't have default as a name", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "notDefault",
				},
			}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).To(HaveOccurred())
			Expect(root).To(BeNil())
		})

		It("should throw when user tries to create non-root singleton object which doesn't have default as a name", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).NotTo(HaveOccurred())
			cfg, err := root.AddConfig(context.TODO(), &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			dns, err := cfg.AddDNS(context.TODO(), &gnsv1.Dns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "notDefault",
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(nexus_client.IsSingletonNameError(err)).To(BeTrue())
			Expect(dns).To(BeNil())
		})

		It("should accept singleton object without a name", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(root.DisplayName()).To(Equal("default"))
			Expect(err).NotTo(HaveOccurred())
			cfg, err := root.AddConfig(context.TODO(), &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			dns, err := cfg.AddDNS(context.TODO(), &gnsv1.Dns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(dns).NotTo(BeNil())
			Expect(dns.DisplayName()).To(Equal("default"))
		})

		It("should accept singleton object without 'default' as a name", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			root, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(root.DisplayName()).To(Equal("default"))
			Expect(err).NotTo(HaveOccurred())
			cfg, err := root.AddConfig(context.TODO(), &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			dns, err := cfg.AddDNS(context.TODO(), &gnsv1.Dns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(dns).NotTo(BeNil())
			Expect(dns.DisplayName()).To(Equal("default"))
		})

		It("shouldn't require singleton object name in chainer", func() {
			fakeClient = nexus_client.NewFakeClient()
			rootDef := &rootv1.Root{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			_, err := fakeClient.AddRootRoot(context.TODO(), rootDef)
			Expect(err).NotTo(HaveOccurred())
			_, err = fakeClient.RootRoot().AddConfig(context.TODO(), &configv1.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			_, err = fakeClient.RootRoot().Config("foo").AddDNS(context.TODO(), &gnsv1.Dns{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
