package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws/autoscaling"
)

func init() {
	RootCmd.AddCommand(newAutoscalingCmd())
}

func newAutoscalingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoscaling",
		Short: "Manage autoscaling resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newAutoscalingLsCmd(),
		newAutoscalingAttachCmd(),
		newAutoscalingDetachCmd(),
		newAutoscalingUpdateCmd(),
	)

	return cmd
}

func newAutoscalingLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List autoscaling groups",
		RunE:  runAutoscalingLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all autoscaling groups (by default, list autoscaling groups only having at least 1 attached instance)")

	viper.BindPFlag("autoscaling.ls.all", flags.Lookup("all"))

	return cmd
}

func runAutoscalingLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := autoscaling.LsOptions{
		All: viper.GetBool("autoscaling.ls.all"),
	}

	return autoscaling.Ls(client, options)
}

func newAutoscalingAttachCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach AUTO_SCALING_GROUP_NAME",
		Short: "Attach instances/loadbalancers to autoscaling group",
		RunE:  runAutoscalingAttachCmd,
	}

	flags := cmd.Flags()
	flags.StringP("instance-ids", "i", "", "One or more instance IDs")
	flags.StringP("load-balancer-names", "l", "", "One or more load balancer names")

	viper.BindPFlag("autoscaling.attach.instance-ids", flags.Lookup("instance-ids"))
	viper.BindPFlag("autoscaling.attach.load-balancer-names", flags.Lookup("load-balancer-names"))

	return cmd
}

func runAutoscalingAttachCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("AUTO_SCALING_GROUP_NAME is required")
	}

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.instance-ids"))
	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.attach.load-balancer-names"))
	options := autoscaling.AttachOptions{
		AsgName:           args[0],
		InstanceIds:       instanceIds,
		LoadBalancerNames: loadBalancerNames,
	}

	return autoscaling.Attach(client, options)
}

func newAutoscalingDetachCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach AUTO_SCALING_GROUP_NAME",
		Short: "Detach instances/loadbalancers from autoscaling group",
		RunE:  runAutoscalingDetachCmd,
	}

	flags := cmd.Flags()
	flags.StringP("instance-ids", "i", "", "One or more instance IDs")
	flags.StringP("load-balancer-names", "l", "", "One or more load balancer names")

	viper.BindPFlag("autoscaling.detach.instance-ids", flags.Lookup("instance-ids"))
	viper.BindPFlag("autoscaling.detach.load-balancer-names", flags.Lookup("load-balancer-names"))

	return cmd
}

func runAutoscalingDetachCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("AUTO_SCALING_GROUP_NAME is required")
	}

	instanceIds := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.instance-ids"))
	loadBalancerNames := aws.StringSlice(viper.GetStringSlice("autoscaling.detach.load-balancer-names"))
	options := autoscaling.DetachOptions{
		AsgName:           args[0],
		InstanceIds:       instanceIds,
		LoadBalancerNames: loadBalancerNames,
	}

	return autoscaling.Detach(client, options)
}

func newAutoscalingUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update AUTO_SCALING_GROUP_NAME",
		Short: "Update autoscaling group",
		RunE:  autoscaling.Update,
	}

	flags := cmd.Flags()
	flags.Int64P("desired-capacity", "c", -1, "The number of EC2 instances that should be running in the Auto Scaling group.")

	viper.BindPFlag("autoscaling.update.desired-capacity", flags.Lookup("desired-capacity"))

	return cmd
}
