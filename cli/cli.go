package cli

import (
	"fmt"
	"github.com/elankath/kcpcl/api"
	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type MainOpts struct {
	api.CopierConfig
	ObjDir                  string
	KubeSchedulerConfigPath string
}

func setupCommonFlagsToOpts(flagSet *flag.FlagSet, mainOpts *MainOpts) {
	flagSet.StringVarP(&mainOpts.KubeConfigPath, clientcmd.RecommendedConfigPathFlag, "k", os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "kubeconfig path of shoot data plane cluster - defaults to KUBECONFIG env-var")
	//downloadFlags.StringVarP(&mainOpts.ControlKubeConfigPath, "kubeconfig-control", "c", os.Getenv("CONTROL_KUBECONFIG"), "kubeconfig path of shoot control plane (seed kubeconfig) - defaults to CONTROL_KUBECONFIG env-var")
	flagSet.StringVarP(&mainOpts.ObjDir, "obj-dir", "d", "", "Base directory where object YAML's of cluster were downloaded using 'download' sub-command")
	flagSet.IntVarP(&mainOpts.PoolSize, "pool-size", "p", 160, "go-routine pool size") //TODO: solve the connection reset by peer issue when pool size increases
}
func SetupDownloadFlagsToOpts(downloadFlags *flag.FlagSet, mainOpts *MainOpts) {
	setupCommonFlagsToOpts(downloadFlags, mainOpts)
	//downloadFlags.StringVarP(&mainOpts.ControlKubeConfigPath, "kubeconfig-control", "c", os.Getenv("CONTROL_KUBECONFIG"), "kubeconfig path of shoot control plane (seed kubeconfig) - defaults to CONTROL_KUBECONFIG env-var")
	standardUsage := downloadFlags.PrintDefaults
	downloadFlags.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s download <flags> <GVRs>\n", api.ProgramName)
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "<flags>")
		standardUsage()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "<GVRs>: GVRs in format [group/][version/]resource where group and version can be omitted for defaults")
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Examples:")
		_, _ = fmt.Fprintf(os.Stderr, "%s download -k /tmp/mykubeconfig.yaml -d /tmp/myobjdir  pods nodes scheduling.k8s.io/v1/priorityclasses\n", api.ProgramName)
		_, _ = fmt.Fprintln(os.Stderr, "  Generate Viewer KubeConfigPath. See: https://github.com/gardener/gardener/blob/23bf7c2dd2e63b338accc68c5b53c1209e9df79a/docs/usage/shoot/shoot_access.md#shootsviewerkubeconfig-subresource")
	}
}
func SetupUploadFlagsToOpts(uploadFlags *flag.FlagSet, mainOpts *MainOpts) {
	setupCommonFlagsToOpts(uploadFlags, mainOpts)
	uploadFlags.StringVarP(&mainOpts.KubeSchedulerConfigPath, "scheduler-config", "s", "/tmp/kube-scheduler-config.yaml", "kube-scheduler config path")
	uploadFlags.BoolVarP(&mainOpts.OrderKinds, "order-kinds", "o", true, "whether to order kinds by priority and wait while uploading")
	standardUsage := uploadFlags.PrintDefaults
	uploadFlags.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s upload <flags>\n", api.ProgramName)
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "<flags>")
		standardUsage()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Examples:")
		_, _ = fmt.Fprintln(os.Stderr, "kcpcl upload -k /tmp/mykubeconfig.yaml -d /tmp/myobjdir")
	}
}

func ValidateMainOptsCommon(mo *MainOpts) (exitCode int, err error) {
	if mo.KubeConfigPath == "" {
		exitCode = ExitMandatoryOpt
		err = api.ErrMissingShootKubeConfig
	}
	if mo.ObjDir == "" {
		exitCode = ExitMandatoryOpt
		err = api.ErrObjDirNotExist
	}
	return
}

func ValidateMainOptsForDownload(mo *MainOpts, args []string) (exitCode int, err error) {
	exitCode, err = ValidateMainOptsCommon(mo)
	if err != nil {
		return
	}
	return
}
func ValidateMainOptsForUpload(mo *MainOpts) (exitCode int, err error) {
	exitCode, err = ValidateMainOptsCommon(mo)
	if err != nil {
		return
	}
	if mo.KubeConfigPath == "" {
		exitCode = ExitMandatoryOpt
		err = api.ErrMissingShootKubeConfig
	}
	if mo.ObjDir == "" {
		exitCode = ExitMandatoryOpt
		err = api.ErrMissingObjDir
	}

	var osFS = afero.NewOsFs()
	ok, err := afero.DirExists(osFS, mo.ObjDir)
	if err != nil {
		exitCode = ExitObjDir
		err = fmt.Errorf("%w: %w", api.ErrCantReadObjDir, err)
		return
	}
	if !ok {
		exitCode = ExitObjDir
		err = fmt.Errorf("%w: %q", api.ErrObjDirNotExist, mo.ObjDir)
		return
	}
	return
}
