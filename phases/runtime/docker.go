package runtime

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/sync/errgroup"
	cmdtext "github.com/yuyicai/kubei/cmd/text"
	"github.com/yuyicai/kubei/config/constants"
	"github.com/yuyicai/kubei/config/rundata"
	"github.com/yuyicai/kubei/phases/system"
	"k8s.io/klog"
)

func InstallDocker(nodes []*rundata.Node, d rundata.Docker) error {
	g := errgroup.WithCancel(context.Background())
	g.GOMAXPROCS(constants.DefaultGOMAXPROCS)
	for _, node := range nodes {
		node := node
		g.Go(func(ctx context.Context) error {
			klog.Infof("[%s] [container-engine] Installing Docker", node.HostInfo.Host)
			if err := installDocker(node, d); err != nil {
				return fmt.Errorf("[%s] [container-engine] Failed to install Docker: %v", node.HostInfo.Host, err)
			}

			if err := system.Restart("docker", node); err != nil {
				return err
			}

			klog.Infof("[%s] [container-engine] Successfully installed Docker", node.HostInfo.Host)
			return nil
		})
	}

	return g.Wait()
}

func installDocker(node *rundata.Node, d rundata.Docker) error {
	cmdText := cmdtext.NewContainerEngineText(node.PackageManagementType)
	cmd, err := cmdText.Docker(node.InstallType, d)
	if err != nil {
		return err
	}

	return node.SSH.Run(cmd)
}
