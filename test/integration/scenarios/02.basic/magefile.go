//go:build mage

/*

	MIT License

	Copyright (c) Microsoft Corporation.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE

*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/princjef/mageutil/shellcmd"
)

// Test config
const (
	TEST_NAME    = "basic manifest deploy scenario"
	TEST_TIMEOUT = "10m"
)

var (
	NAMESPACES = []string{
		"nondefault",
		"default",
	}
)

var (
	// Manifests to deploy
	testManifests = []string{
		"manifest/%s/target.yaml",
		"manifest/%s/instance.yaml",
		"manifest/%s/solution.yaml",
	}

	// Tests to run
	testVerify = []string{
		"./verify/...",
	}
)

// Entry point for running the tests
func Test() error {
	fmt.Println("Running ", TEST_NAME)

	defer Cleanup()
	err := SetupCluster()
	if err != nil {
		return err
	}
	for _, namespace := range NAMESPACES {
		os.Setenv("NAMESPACE", namespace)
		err = DeployManifests(namespace)
		if err != nil {
			return err
		}
		err = Verify()
		if err != nil {
			return err
		}
		err = CleanUpSymphonyObjects(namespace)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 10)
	}

	return nil
}

// Prepare the cluster
// Run this manually to prepare your local environment for testing/debugging
func SetupCluster() error {
	// Deploy symphony
	err := localenvCmd("cluster:deploy", "")
	if err != nil {
		return err
	}
	// Wait a few secs for symphony cert to be ready;
	// otherwise we will see error when creating symphony manifests in the cluster
	// <Error from server (InternalError): error when creating
	// "/mnt/vss/_work/1/s/test/integration/scenarios/basic/manifest/target.yaml":
	// Internal error occurred: failed calling webhook "mtarget.kb.io": failed to
	// call webhook: Post
	// "https://symphony-webhook-service.default.svc:443/mutate-symphony-microsoft-com-v1-target?timeout=10s":
	// x509: certificate signed by unknown authority>
	time.Sleep(time.Second * 10)
	return nil
}

func DeployManifests(namespace string) error {
	if namespace != "default" {
		// Create non-default namespace if not exist
		err := shellcmd.Command(fmt.Sprintf("kubectl get namespace %s", namespace)).Run()
		if err != nil {
			// Better to check err message here but command only returns "exit status 1" for non-exisiting namespace
			err = shellcmd.Command(fmt.Sprintf("kubectl create namespace %s", namespace)).Run()
			if err != nil {
				return err
			}
		}
	}
	// Deploy the manifests
	for _, manifest := range testManifests {
		fullPath, err := filepath.Abs(fmt.Sprintf(manifest, "oss"))
		if err != nil {
			return err
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}
		stringYaml := string(data)
		stringYaml = strings.ReplaceAll(stringYaml, "INSTANCENAME", namespace+"instance")
		stringYaml = strings.ReplaceAll(stringYaml, "SCOPENAME", namespace+"scope")
		stringYaml = strings.ReplaceAll(stringYaml, "TARGETNAME", namespace+"target")
		stringYaml = strings.ReplaceAll(stringYaml, "SOLUTIONNAME", namespace+"solution")

		err = writeYamlStringsToFile(stringYaml, "./test.yaml")
		if err != nil {
			return err
		}
		err = shellcmd.Command(fmt.Sprintf("kubectl apply -f ./test.yaml -n %s", namespace)).Run()
		if err != nil {
			return err
		}
		os.Remove("./test.yaml")
	}

	return nil
}

// Run tests
func Verify() error {
	err := shellcmd.Command("go clean -testcache").Run()
	if err != nil {
		return err
	}
	os.Setenv("SYMPHONY_FLAVOR", "oss")
	for _, verify := range testVerify {
		err := shellcmd.Command(fmt.Sprintf("go test -timeout %s %s", TEST_TIMEOUT, verify)).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanUpSymphonyObjects(namespace string) error {
	instanceName := namespace + "instance"
	targetName := namespace + "target"
	solutionName := namespace + "solution"
	err := shellcmd.Command(fmt.Sprintf("kubectl delete instance %s -n %s", instanceName, namespace)).Run()
	if err != nil {
		return err
	}
	err = shellcmd.Command(fmt.Sprintf("kubectl delete target %s -n %s", targetName, namespace)).Run()
	if err != nil {
		return err
	}
	err = shellcmd.Command(fmt.Sprintf("kubectl delete solution %s -n %s", solutionName, namespace)).Run()
	if err != nil {
		return err
	}
	return nil
}

// Clean up
func Cleanup() {
	localenvCmd("destroy all", "")
}

// Run a mage command from /localenv
func localenvCmd(mageCmd string, flavor string) error {
	return shellExec(fmt.Sprintf("cd ../../../localenv && mage %s %s", mageCmd, flavor))
}

func shellExec(cmd string) error {
	fmt.Println("> ", cmd)

	execCmd := exec.Command("sh", "-c", cmd)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

func writeYamlStringsToFile(yamlString string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(yamlString))
	if err != nil {
		return err
	}

	return nil
}
