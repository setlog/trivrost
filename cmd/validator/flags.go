package main

import "flag"

type ValidatorFlags struct {
	DeploymentConfigUrl string
	SkipUrlCheck        bool
	SkipJarChek         bool
	ActAsService        bool
	Port                int
}

func parseFlags() *ValidatorFlags {
	flags := ValidatorFlags{}
	flag.BoolVar(&flags.SkipUrlCheck, "skipurlcheck", false, "Disable checking of availability of all URLs in the config.")
	flag.BoolVar(&flags.SkipJarChek, "skipjarcheck", false, "Disable checking of availability of .jar files given to java with the -jar argument.")
	flag.BoolVar(&flags.ActAsService, "act-as-service", false, "Validate deployment-config for HTTP GET requests on :80/validate.")
	flag.IntVar(&flags.Port, "port", 80, "Override port for --act-as-service.")
	flag.Parse()

	if flag.NArg() > 1 {
		fatalf("Too many arguments. Required: deploymentConfigURL")
	} else if flag.NArg() == 1 {
		flags.DeploymentConfigUrl = flag.Arg(0)
	} else if !flags.ActAsService {
		fatalf("The following argument is required when not running with --act-as-service: deploymentConfigURL")
	}

	return &flags
}
