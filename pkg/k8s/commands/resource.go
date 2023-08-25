package commands

import (
	"context"
	"strings"

	"github.com/khulnasoft-lab/vul/pkg/flag"

	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/vul-kubernetes/pkg/artifacts"
	"github.com/khulnasoft-lab/vul-kubernetes/pkg/k8s"
	"github.com/khulnasoft-lab/vul-kubernetes/pkg/vulk8s"
	"github.com/khulnasoft-lab/vul/pkg/log"
)

// resourceRun runs scan on kubernetes cluster
func resourceRun(ctx context.Context, args []string, opts flag.Options, cluster k8s.Cluster) error {
	kind, name, err := extractKindAndName(args)
	if err != nil {
		return err
	}

	vulk8s := vulk8s.New(cluster, log.Logger).Namespace(getNamespace(opts, cluster.GetCurrentNamespace()))
	runner := newRunner(opts, cluster.GetCurrentContext())

	if len(name) == 0 { // pods or configmaps etc
		if err = validateReportArguments(opts); err != nil {
			return err
		}

		targets, err := vulk8s.Resources(kind).ListArtifacts(ctx)
		if err != nil {
			return err
		}

		return runner.run(ctx, targets)
	}

	// pod/NAME or pod NAME etc
	artifact, err := vulk8s.GetArtifact(ctx, kind, name)
	if err != nil {
		return err
	}

	return runner.run(ctx, []*artifacts.Artifact{artifact})
}

func extractKindAndName(args []string) (string, string, error) {
	switch len(args) {
	case 1:
		s := strings.Split(args[0], "/")
		if len(s) != 2 {
			return args[0], "", nil
		}

		return s[0], s[1], nil
	case 2:
		return args[0], args[1], nil
	}

	return "", "", xerrors.Errorf("can't parse arguments %v. Please run `vul k8s` for usage.", args)
}
