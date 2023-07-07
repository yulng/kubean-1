package kubean_non_root_user_create_cluster

import (
	"fmt"
	"strings"

	"github.com/kubean-io/kubean/test/tools"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var _ = ginkgo.Describe("non root users create a work cluster", func() {

	ginkgo.Context("Parameters init", func() {
		var masterSSH = fmt.Sprintf("root@%s", tools.Vmipaddr)
		var master2SSH = fmt.Sprintf("root@%s", tools.Vmipaddr2)
		var master3SSH = fmt.Sprintf("root@%s", tools.Vmipaddr3)
		var vipaddcr = tools.Vipadd
		localKubeConfigPath := "cluster1.config"
		var offlineConfigs tools.OfflineConfig
		var password = tools.VmPassword
		testClusterName := tools.TestClusterName
		nginxImage := tools.NginxAlpha
		offlineFlag := tools.IsOffline
		offlineConfigs = tools.InitOfflineConfig()
		if strings.ToUpper(offlineFlag) == "TRUE" && strings.ToUpper(tools.Arch) == "ARM64" {
			nginxImage = offlineConfigs.NginxImageARM64
		}
		if strings.ToUpper(offlineFlag) == "TRUE" && strings.ToUpper(tools.Arch) == "AMD64" {
			nginxImage = offlineConfigs.NginxImageAMD64
		}
		klog.Info("nginx image is: ", nginxImage)
		klog.Info("offlineFlag is: ", offlineFlag)
		klog.Info("arch is: ", tools.Arch)

		// non root users create a work cluster
		ginkgo.It("non root users create a work cluster", func() {

			clusterInstallYamlsPath := tools.E2eInstallClusterYamlFolder
			kubeanClusterOpsName := tools.ClusterOperationName
			kindConfig, err := clientcmd.BuildConfigFromFlags("", tools.Kubeconfig)
			gomega.ExpectWithOffset(2, err).NotTo(gomega.HaveOccurred(), "failed build config")
			tools.OperateClusterByYaml(clusterInstallYamlsPath, kubeanClusterOpsName, kindConfig)

			tools.SaveKubeConf(kindConfig, testClusterName, localKubeConfigPath)
			cluster1Config, err := clientcmd.BuildConfigFromFlags("", localKubeConfigPath)
			gomega.ExpectWithOffset(2, err).NotTo(gomega.HaveOccurred(), "Failed new cluster1Config set")
			cluster1Client, err := kubernetes.NewForConfig(cluster1Config)
			gomega.ExpectWithOffset(2, err).NotTo(gomega.HaveOccurred(), "Failed new cluster1Client")
			tools.WaitPodSInKubeSystemBeRunning(cluster1Client, 1800)
			// do sonobuoy check
			if strings.ToUpper(offlineFlag) != "TRUE" {
				klog.Info("On line, sonobuoy check")
				tools.DoSonoBuoyCheckByPasswd(password, masterSSH)
			}
		})

	})
})
