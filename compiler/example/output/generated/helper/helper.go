package helper

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/elliotchance/orderedmap"

	datamodel "nexustempmodule/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DEFAULT_KEY = "default"
const DISPLAY_NAME_LABEL = "nexus/display_name"
const IS_NAME_HASHED_LABEL = "nexus/is_name_hashed"

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"accesscontrolpolicies.global.tsm.tanzu.vmware.com":                  {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"accesstokens.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "userfolders.global.tsm.tanzu.vmware.com"},
		"acpconfigs.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "accesscontrolpolicies.global.tsm.tanzu.vmware.com"},
		"additionalattributeses.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com", "services.global.tsm.tanzu.vmware.com"},
		"allsparkserviceses.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"annotations.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "globalregistrationservices.global.tsm.tanzu.vmware.com", "tenants.global.tsm.tanzu.vmware.com"},
		"apidiscoveries.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"apidiscoveryrts.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"appfolders.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"appgroups.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "resourcegroups.global.tsm.tanzu.vmware.com"},
		"applicationinfos.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"apps.global.tsm.tanzu.vmware.com":                                   {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "appfolders.global.tsm.tanzu.vmware.com"},
		"apptemplates.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com"},
		"apptemplateservicedefinitions.global.tsm.tanzu.vmware.com":          {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com", "apptemplates.global.tsm.tanzu.vmware.com"},
		"appusers.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "userfolders.global.tsm.tanzu.vmware.com"},
		"appversions.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "appfolders.global.tsm.tanzu.vmware.com", "apps.global.tsm.tanzu.vmware.com"},
		"attackdiscoveries.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"attackdiscoveryrts.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"authenticationpolicies.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"autoscalerconfigs.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"autoscalercrs.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"autoscalerfolders.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"autoscalers.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "autoscalerfolders.global.tsm.tanzu.vmware.com"},
		"autoscalingpolicies.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"awsconnectors.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com"},
		"buckets.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com", "databases.global.tsm.tanzu.vmware.com"},
		"certificateauthorityconfigns.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "projects.global.tsm.tanzu.vmware.com", "projectconfigs.global.tsm.tanzu.vmware.com"},
		"certificateauthorityrts.global.tsm.tanzu.vmware.com":                {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"certificateconfigns.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"certificaterequests.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"certificates.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"clusterconfigfolders.global.tsm.tanzu.vmware.com":                   {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"clusterconfigs.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com"},
		"clusterfolders.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"clusters.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com"},
		"clustersettingses.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "clusterconfigfolders.global.tsm.tanzu.vmware.com"},
		"configmaps.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"configs.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com"},
		"connectionstatuses.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"databases.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com"},
		"datafolderdomainclusters.global.tsm.tanzu.vmware.com":               {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com"},
		"datafolderdomains.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com"},
		"datafolderdomainservices.global.tsm.tanzu.vmware.com":               {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com"},
		"datafolderdomainserviceversions.global.tsm.tanzu.vmware.com":        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com", "datafolderdomainservices.global.tsm.tanzu.vmware.com"},
		"datafolders.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"datagroups.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "resourcegroups.global.tsm.tanzu.vmware.com"},
		"datatemplates.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com"},
		"dcregions.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com"},
		"dczones.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "dcregions.global.tsm.tanzu.vmware.com"},
		"destinationrules.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"directories.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com", "databases.global.tsm.tanzu.vmware.com"},
		"dnsconfigfolders.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"dnsconfigs.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "dnsconfigfolders.global.tsm.tanzu.vmware.com"},
		"dnsprobeconfigs.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "dnsprobesconfigfolders.global.tsm.tanzu.vmware.com"},
		"dnsprobesconfigfolders.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"dnsprobestatuses.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"domainconfigs.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"domains.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"endpoints.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"envoyfilters.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"externalaccountconfigns.global.tsm.tanzu.vmware.com":                {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"externalauditstorages.global.tsm.tanzu.vmware.com":                  {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"externaldnsconfigns.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"externaldnsinventories.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com"},
		"externaldnsinventoryhealthchecks.global.tsm.tanzu.vmware.com":       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "externaldnsinventories.global.tsm.tanzu.vmware.com"},
		"externaldnsinventoryprimarydomains.global.tsm.tanzu.vmware.com":     {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "externaldnsinventories.global.tsm.tanzu.vmware.com"},
		"externaldnsinventoryrecords.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "externaldnsinventories.global.tsm.tanzu.vmware.com", "externaldnsinventoryprimarydomains.global.tsm.tanzu.vmware.com"},
		"externaldnsinventoryzones.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "externaldnsinventories.global.tsm.tanzu.vmware.com"},
		"externaldnsruntimeendpoints.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "externaldnsruntimes.global.tsm.tanzu.vmware.com", "externaldnsruntimeprimarydomains.global.tsm.tanzu.vmware.com", "externaldnsruntimesubdomains.global.tsm.tanzu.vmware.com"},
		"externaldnsruntimeprimarydomains.global.tsm.tanzu.vmware.com":       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "externaldnsruntimes.global.tsm.tanzu.vmware.com"},
		"externaldnsruntimes.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"externaldnsruntimesubdomains.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "externaldnsruntimes.global.tsm.tanzu.vmware.com", "externaldnsruntimeprimarydomains.global.tsm.tanzu.vmware.com"},
		"externallbconfigns.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"externalplugincapabilities.global.tsm.tanzu.vmware.com":             {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "externalpluginconfigfolders.global.tsm.tanzu.vmware.com", "externalpluginconfigs.global.tsm.tanzu.vmware.com", "externalplugininstanceconfigs.global.tsm.tanzu.vmware.com"},
		"externalpluginconfigfolders.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"externalpluginconfigs.global.tsm.tanzu.vmware.com":                  {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "externalpluginconfigfolders.global.tsm.tanzu.vmware.com"},
		"externalplugininstanceconfigs.global.tsm.tanzu.vmware.com":          {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "externalpluginconfigfolders.global.tsm.tanzu.vmware.com", "externalpluginconfigs.global.tsm.tanzu.vmware.com"},
		"externalpluginmonitors.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "externalpluginconfigfolders.global.tsm.tanzu.vmware.com", "externalpluginconfigs.global.tsm.tanzu.vmware.com", "externalplugininstanceconfigs.global.tsm.tanzu.vmware.com"},
		"externalserviceconfigs.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"externalservicesrts.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"featureflags.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"federatedsloconfigs.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "slofolders.global.tsm.tanzu.vmware.com"},
		"federatedsloserviceconfigs.global.tsm.tanzu.vmware.com":             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "slofolders.global.tsm.tanzu.vmware.com", "federatedsloconfigs.global.tsm.tanzu.vmware.com"},
		"gatewayconfigadditionallistenerses.global.tsm.tanzu.vmware.com":     {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "gatewayconfigs.global.tsm.tanzu.vmware.com"},
		"gatewayconfiglistenercertificates.global.tsm.tanzu.vmware.com":      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "gatewayconfigs.global.tsm.tanzu.vmware.com", "gatewayconfigadditionallistenerses.global.tsm.tanzu.vmware.com"},
		"gatewayconfiglistenerdestinationroutes.global.tsm.tanzu.vmware.com": {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "gatewayconfigs.global.tsm.tanzu.vmware.com", "gatewayconfigadditionallistenerses.global.tsm.tanzu.vmware.com"},
		"gatewayconfigs.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"gateways.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"geodiscoveries.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"geodiscoveryrts.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"globalnamespaces.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"globalnses.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "awsconnectors.global.tsm.tanzu.vmware.com"},
		"globalregistrationservices.global.tsm.tanzu.vmware.com":             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com"},
		"gnsaccesscontrolpolicies.global.tsm.tanzu.vmware.com":               {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"gnsaccesscontrolpolicyrts.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"gnsbindingrts.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"gnsendpointsconfigs.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"gnsroutingconfigs.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"gnsroutingruleconfigs.global.tsm.tanzu.vmware.com":                  {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "gnsbindingrts.global.tsm.tanzu.vmware.com"},
		"gnss.global.tsm.tanzu.vmware.com":                                   {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com"},
		"gnsschemaviolationdiscoveries.global.tsm.tanzu.vmware.com":          {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"gnssegmentationpolicies.global.tsm.tanzu.vmware.com":                {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"gnssegmentationpolicyrts.global.tsm.tanzu.vmware.com":               {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"gnsserviceentryconfigs.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "gnsbindingrts.global.tsm.tanzu.vmware.com"},
		"gnssvcgrouprts.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "gnsbindingrts.global.tsm.tanzu.vmware.com"},
		"gnssvcgroups.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"haconfigs.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "hostconfigs.global.tsm.tanzu.vmware.com"},
		"haconfigv2s.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "hostconfigv2s.global.tsm.tanzu.vmware.com"},
		"healthcheckconfigns.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com"},
		"hostconfigs.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"hostconfigv2s.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"inboundauthenticationconfigs.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"inventories.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com"},
		"issuers.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"jobconfigfolders.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"jobconfigs.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "jobconfigfolders.global.tsm.tanzu.vmware.com"},
		"jobfolders.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com"},
		"jobs.global.tsm.tanzu.vmware.com":                                   {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "jobfolders.global.tsm.tanzu.vmware.com"},
		"knativeingresses.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"labelconfigs.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"localregistrationserviceclusters.global.tsm.tanzu.vmware.com":       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "localregistrationservices.global.tsm.tanzu.vmware.com"},
		"localregistrationserviceresources.global.tsm.tanzu.vmware.com":      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "localregistrationservices.global.tsm.tanzu.vmware.com", "localregistrationserviceclusters.global.tsm.tanzu.vmware.com"},
		"localregistrationservices.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com"},
		"logfolders.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"logs.global.tsm.tanzu.vmware.com":                                   {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "logfolders.global.tsm.tanzu.vmware.com"},
		"metricmonitors.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "serviceconfigs.global.tsm.tanzu.vmware.com", "serviceversionconfigs.global.tsm.tanzu.vmware.com"},
		"networkattachmentdefinitionconfigs.global.tsm.tanzu.vmware.com":     {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"networkattachmentdefinitions.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"nodedefinitions.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com", "nodetemplates.global.tsm.tanzu.vmware.com"},
		"nodefolderclusters.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "nodefolders.global.tsm.tanzu.vmware.com"},
		"nodefolders.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"nodegroups.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "resourcegroups.global.tsm.tanzu.vmware.com"},
		"nodes.global.tsm.tanzu.vmware.com":                                  {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"nodestatuses.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "nodes.global.tsm.tanzu.vmware.com"},
		"nodetemplates.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com"},
		"outboundauthenticationmodes.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "hostconfigv2s.global.tsm.tanzu.vmware.com"},
		"peerauthentications.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"piidiscoveries.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"piidiscoveryrts.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"policyconfigs.global.tsm.tanzu.vmware.com":                          {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"policytemplates.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com"},
		"progressiveupgradeconfigs.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "progressiveupgrades.global.tsm.tanzu.vmware.com"},
		"progressiveupgradefolders.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"progressiveupgraderuntimes.global.tsm.tanzu.vmware.com":             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "progressiveupgradefolders.global.tsm.tanzu.vmware.com"},
		"progressiveupgrades.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"projectconfigs.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "projects.global.tsm.tanzu.vmware.com"},
		"projectinventories.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "projects.global.tsm.tanzu.vmware.com"},
		"projectqueries.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "projects.global.tsm.tanzu.vmware.com"},
		"projects.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"publicserviceconfigs.global.tsm.tanzu.vmware.com":                   {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"publicservicerouteconfigs.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com", "publicserviceconfigs.global.tsm.tanzu.vmware.com"},
		"publicservicerts.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"remotegatewayserviceconfigs.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com"},
		"resourcegrouprts.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"resourcegroups.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"roots.global.tsm.tanzu.vmware.com":                                  {},
		"rpolicies.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "awsconnectors.global.tsm.tanzu.vmware.com", "globalnses.global.tsm.tanzu.vmware.com"},
		"runtimes.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com"},
		"schemaviolationdiscoveryrts.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"secrethashes.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"secretrtconfigs.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"securitycontextconstraintsconfigs.global.tsm.tanzu.vmware.com":      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com"},
		"securitycontextconstraintses.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com"},
		"serviceconfigs.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"servicecronjobs.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicedaemonsets.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicedeploymentcontainers.global.tsm.tanzu.vmware.com":            {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com", "servicedeployments.global.tsm.tanzu.vmware.com"},
		"servicedeployments.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicedirectoryentryconfigs.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "servicedirectoryns.global.tsm.tanzu.vmware.com"},
		"servicedirectoryns.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"servicedirectoryrtfolderentries.global.tsm.tanzu.vmware.com":        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "servicedirectoryrts.global.tsm.tanzu.vmware.com", "servicedirectoryrtfolders.global.tsm.tanzu.vmware.com"},
		"servicedirectoryrtfolders.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "servicedirectoryrts.global.tsm.tanzu.vmware.com"},
		"servicedirectoryrts.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"serviceentries.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"serviceentryconfigs.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"serviceinstancecontainers.global.tsm.tanzu.vmware.com":              {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com", "serviceinstances.global.tsm.tanzu.vmware.com"},
		"serviceinstances.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicejobs.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicelevelobjectivefolders.global.tsm.tanzu.vmware.com":           {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"servicelevelobjectives.global.tsm.tanzu.vmware.com":                 {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "servicelevelobjectivefolders.global.tsm.tanzu.vmware.com"},
		"servicereplicasets.global.tsm.tanzu.vmware.com":                     {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"services.global.tsm.tanzu.vmware.com":                               {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicestatefulsets.global.tsm.tanzu.vmware.com":                    {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"servicetemplates.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com"},
		"servicetemplateservicedefinitions.global.tsm.tanzu.vmware.com":      {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com", "templategroups.global.tsm.tanzu.vmware.com", "servicetemplates.global.tsm.tanzu.vmware.com"},
		"serviceversionconfigs.global.tsm.tanzu.vmware.com":                  {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "serviceconfigs.global.tsm.tanzu.vmware.com"},
		"serviceversions.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com", "services.global.tsm.tanzu.vmware.com"},
		"sharedserviceconfigs.global.tsm.tanzu.vmware.com":                   {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"sloconfigs.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"slofolders.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"slopolicies.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"sloserviceconfigs.global.tsm.tanzu.vmware.com":                      {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com", "sloconfigs.global.tsm.tanzu.vmware.com"},
		"svcgrouprts.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "resourcegrouprts.global.tsm.tanzu.vmware.com"},
		"svcgroups.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "resourcegroups.global.tsm.tanzu.vmware.com"},
		"tables.global.tsm.tanzu.vmware.com":                                 {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "datafolders.global.tsm.tanzu.vmware.com", "datafolderdomains.global.tsm.tanzu.vmware.com", "databases.global.tsm.tanzu.vmware.com"},
		"templategroups.global.tsm.tanzu.vmware.com":                         {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "templates.global.tsm.tanzu.vmware.com"},
		"templates.global.tsm.tanzu.vmware.com":                              {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com"},
		"tenantresources.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "globalregistrationservices.global.tsm.tanzu.vmware.com", "tenants.global.tsm.tanzu.vmware.com"},
		"tenants.global.tsm.tanzu.vmware.com":                                {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "globalregistrationservices.global.tsm.tanzu.vmware.com"},
		"tenanttokens.global.tsm.tanzu.vmware.com":                           {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "allsparkserviceses.global.tsm.tanzu.vmware.com", "globalregistrationservices.global.tsm.tanzu.vmware.com"},
		"userdiscoveries.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "globalnamespaces.global.tsm.tanzu.vmware.com", "gnss.global.tsm.tanzu.vmware.com"},
		"userdiscoveryrts.global.tsm.tanzu.vmware.com":                       {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "clusterfolders.global.tsm.tanzu.vmware.com", "clusterconfigs.global.tsm.tanzu.vmware.com", "domainconfigs.global.tsm.tanzu.vmware.com"},
		"userfolders.global.tsm.tanzu.vmware.com":                            {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com"},
		"usergroups.global.tsm.tanzu.vmware.com":                             {"roots.global.tsm.tanzu.vmware.com", "configs.global.tsm.tanzu.vmware.com", "resourcegroups.global.tsm.tanzu.vmware.com"},
		"userpreferences.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "userfolders.global.tsm.tanzu.vmware.com", "users.global.tsm.tanzu.vmware.com"},
		"users.global.tsm.tanzu.vmware.com":                                  {"roots.global.tsm.tanzu.vmware.com", "runtimes.global.tsm.tanzu.vmware.com", "userfolders.global.tsm.tanzu.vmware.com"},
		"virtualservices.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
		"workloadentries.global.tsm.tanzu.vmware.com":                        {"roots.global.tsm.tanzu.vmware.com", "inventories.global.tsm.tanzu.vmware.com", "clusters.global.tsm.tanzu.vmware.com", "domains.global.tsm.tanzu.vmware.com"},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "accesscontrolpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AccessControlPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "accesstokens.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AccessTokens().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "acpconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ACPConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "additionalattributeses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AdditionalAttributeses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "allsparkserviceses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AllSparkServiceses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "annotations.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Annotations().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "apidiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ApiDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "apidiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ApiDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "appfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "appgroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "applicationinfos.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ApplicationInfos().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "apps.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Apps().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "apptemplates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "apptemplateservicedefinitions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppTemplateServiceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "appusers.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppUsers().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "appversions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AppVersions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "attackdiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AttackDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "attackdiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AttackDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "authenticationpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AuthenticationPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "autoscalerconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AutoscalerConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "autoscalercrs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AutoscalerCRs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "autoscalerfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AutoscalerFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "autoscalers.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Autoscalers().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "autoscalingpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AutoscalingPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "awsconnectors.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().AwsConnectors().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "buckets.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Buckets().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "certificateauthorityconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().CertificateAuthorityConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "certificateauthorityrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().CertificateAuthorityRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "certificateconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().CertificateConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "certificaterequests.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().CertificateRequests().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "certificates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Certificates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "clusterconfigfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ClusterConfigFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "clusterconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ClusterConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "clusterfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ClusterFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "clusters.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Clusters().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "clustersettingses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ClusterSettingses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "configmaps.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ConfigMaps().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "configs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Configs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "connectionstatuses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ConnectionStatuses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "databases.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Databases().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datafolderdomainclusters.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataFolderDomainClusters().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datafolderdomains.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataFolderDomains().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datafolderdomainservices.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataFolderDomainServices().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datafolderdomainserviceversions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataFolderDomainServiceVersions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datafolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datagroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "datatemplates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DataTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dcregions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DCRegions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dczones.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DCZones().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "destinationrules.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DestinationRules().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "directories.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Directories().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dnsconfigfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DNSConfigFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dnsconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DNSConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dnsprobeconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DNSProbeConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dnsprobesconfigfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DNSProbesConfigFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "dnsprobestatuses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DNSProbeStatuses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "domainconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().DomainConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "domains.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Domains().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "endpoints.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Endpoints().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "envoyfilters.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().EnvoyFilters().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalaccountconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalAccountConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalauditstorages.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalAuditStorages().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsinventories.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSInventories().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsinventoryhealthchecks.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSInventoryHealthChecks().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsinventoryprimarydomains.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSInventoryPrimaryDomains().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsinventoryrecords.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSInventoryRecords().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsinventoryzones.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSInventoryZones().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsruntimeendpoints.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSRuntimeEndpoints().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsruntimeprimarydomains.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSRuntimePrimaryDomains().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsruntimes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSRuntimes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externaldnsruntimesubdomains.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalDNSRuntimeSubdomains().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externallbconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalLBConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalplugincapabilities.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalPluginCapabilities().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalpluginconfigfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalPluginConfigFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalpluginconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalPluginConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalplugininstanceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalPluginInstanceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalpluginmonitors.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalPluginMonitors().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "externalservicesrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ExternalServicesRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "featureflags.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().FeatureFlags().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "federatedsloconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().FederatedSloConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "federatedsloserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().FederatedSloServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gatewayconfigadditionallistenerses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GatewayConfigAdditionalListenerses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gatewayconfiglistenercertificates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GatewayConfigListenerCertificates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gatewayconfiglistenerdestinationroutes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GatewayConfigListenerDestinationRoutes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gatewayconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GatewayConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gateways.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Gateways().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "geodiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GeoDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "geodiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GeoDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "globalnamespaces.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GlobalNamespaces().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "globalnses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GlobalNses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "globalregistrationservices.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GlobalRegistrationServices().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsaccesscontrolpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsAccessControlPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsaccesscontrolpolicyrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsAccessControlPolicyRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsbindingrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsBindingRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsendpointsconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsEndpointsConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsroutingconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GNSRoutingConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsroutingruleconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsRoutingRuleConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnss.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GNSs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsschemaviolationdiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsSchemaViolationDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnssegmentationpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsSegmentationPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnssegmentationpolicyrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsSegmentationPolicyRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnsserviceentryconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsServiceEntryConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnssvcgrouprts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GnsSvcGroupRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnssvcgroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().GNSSvcGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "haconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().HaConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "haconfigv2s.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().HaConfigV2s().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "healthcheckconfigns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().HealthCheckConfigNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "hostconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().HostConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "hostconfigv2s.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().HostConfigV2s().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "inboundauthenticationconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().InboundAuthenticationConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "inventories.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Inventories().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "issuers.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Issuers().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "jobconfigfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().JobConfigFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "jobconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().JobConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "jobfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().JobFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "jobs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Jobs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "knativeingresses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().KnativeIngresses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "labelconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().LabelConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "localregistrationserviceclusters.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().LocalRegistrationServiceClusters().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "localregistrationserviceresources.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().LocalRegistrationServiceResources().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "localregistrationservices.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().LocalRegistrationServices().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "logfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().LogFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "logs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Logs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "metricmonitors.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().MetricMonitors().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "networkattachmentdefinitionconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NetworkAttachmentDefinitionConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "networkattachmentdefinitions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NetworkAttachmentDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodedefinitions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodefolderclusters.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeFolderClusters().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodefolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodegroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodestatuses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeStatuses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nodetemplates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().NodeTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "outboundauthenticationmodes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().OutboundAuthenticationModes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "peerauthentications.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PeerAuthentications().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "piidiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PiiDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "piidiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PiiDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "policyconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PolicyConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "policytemplates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PolicyTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "progressiveupgradeconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProgressiveUpgradeConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "progressiveupgradefolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProgressiveUpgradeFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "progressiveupgraderuntimes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProgressiveUpgradeRuntimes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "progressiveupgrades.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProgressiveUpgrades().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "projectconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProjectConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "projectinventories.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProjectInventories().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "projectqueries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ProjectQueries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "projects.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Projects().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "publicserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PublicServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "publicservicerouteconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PublicServiceRouteConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "publicservicerts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().PublicServiceRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "remotegatewayserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().RemoteGatewayServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "resourcegrouprts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ResourceGroupRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "resourcegroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ResourceGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "roots.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Roots().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "rpolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().RPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "runtimes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Runtimes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "schemaviolationdiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SchemaViolationDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "secrethashes.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SecretHashes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "secretrtconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SecretRTConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "securitycontextconstraintsconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SecurityContextConstraintsConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "securitycontextconstraintses.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SecurityContextConstraintses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicecronjobs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceCronJobs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedaemonsets.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDaemonSets().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedeploymentcontainers.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDeploymentContainers().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedeployments.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDeployments().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedirectoryentryconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDirectoryEntryConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedirectoryns.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDirectoryNs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedirectoryrtfolderentries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDirectoryRTFolderEntries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedirectoryrtfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDirectoryRTFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicedirectoryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceDirectoryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceentries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceEntries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceentryconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceEntryConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceinstancecontainers.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceInstanceContainers().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceinstances.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceInstances().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicejobs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceJobs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicelevelobjectivefolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceLevelObjectiveFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicelevelobjectives.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceLevelObjectives().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicereplicasets.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceReplicaSets().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "services.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Services().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicestatefulsets.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceStatefulSets().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicetemplates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "servicetemplateservicedefinitions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceTemplateServiceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceversionconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceVersionConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "serviceversions.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().ServiceVersions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "sharedserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SharedServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "sloconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SloConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "slofolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SLOFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "slopolicies.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SLOPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "sloserviceconfigs.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SloServiceConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "svcgrouprts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SvcGroupRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "svcgroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().SvcGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tables.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Tables().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "templategroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().TemplateGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "templates.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Templates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tenantresources.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().TenantResources().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tenants.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Tenants().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tenanttokens.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().TenantTokens().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "userdiscoveries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().UserDiscoveries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "userdiscoveryrts.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().UserDiscoveryRTs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "userfolders.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().UserFolders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "usergroups.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().UserGroups().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "userpreferences.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().UserPreferences().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "users.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().Users().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "virtualservices.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().VirtualServices().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "workloadentries.global.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GlobalTsmV1().WorkloadEntries().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}

	return nil
}

func ParseCRDLabels(crdName string, labels map[string]string) *orderedmap.OrderedMap {
	parents := GetCRDParentsMap()[crdName]

	m := orderedmap.NewOrderedMap()
	for _, parent := range parents {
		if label, ok := labels[parent]; ok {
			m.Set(parent, label)
		} else {
			m.Set(parent, DEFAULT_KEY)
		}
	}

	return m
}

func GetHashedName(crdName string, labels map[string]string, name string) string {
	orderedLabels := ParseCRDLabels(crdName, labels)

	var output string
	for i, key := range orderedLabels.Keys() {
		value, _ := orderedLabels.Get(key)

		output += fmt.Sprintf("%s:%s", key, value)
		if i < orderedLabels.Len()-1 {
			output += "/"
		}
	}

	output += fmt.Sprintf("%s:%s", crdName, name)

	h := sha1.New()
	_, _ = h.Write([]byte(output))
	return hex.EncodeToString(h.Sum(nil))
}
