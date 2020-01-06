package main

import "flag"

type ValidatorFlags struct {
	DeploymentConfigUrl string
	SkipUrlCheck        bool
	SkipJarChek         bool
	ActAsService        bool
	Port                int
}

func parseFlags() ValidatorFlags {
	flags := ValidatorFlags{}
	flag.BoolVar(&flags.SkipUrlCheck, "skipurlcheck", false, "Disable checking of availability of all URLs in the config.")
	flag.BoolVar(&flags.SkipJarChek, "skipjarcheck", false, "Disable checking of availability of .jar files given to java with the -jar argument.")
	flag.BoolVar(&flags.ActAsService, "act-as-service", false, "Validate deployment-config for HTTP GET requests on :80/validate.")
	flag.IntVar(&flags.Port, "port", 80, "Override port for --act-as-service.")
	flag.Parse()

	if flag.NArg() != 1 {
		fatalf("Need at least one arg: deploymentConfigURL")
	}
	flags.DeploymentConfigUrl = flag.Arg(0)
	if flags.DeploymentConfigUrl == "" {
		fatalf("deploymentConfigURL not set")
	}

	return flags
}
