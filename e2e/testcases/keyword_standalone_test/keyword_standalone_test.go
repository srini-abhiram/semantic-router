package testcases

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/vllm-project/semantic-router/e2e/pkg/testcases"
	"github.com/vllm-project/semantic-router/src/semantic-router/pkg/classification"
	"github.com/vllm-project/semantic-router/src/semantic-router/pkg/config"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
)

func init() {
	testcases.Register("keyword-standalone-test", testcases.TestCase{
		Description: "Test keyword-based routing in a standalone manner",
		Tags:        []string{"keyword", "standalone"},
		Fn:          testKeywordStandalone,
	})
}

type IntelligentRoute struct {
	Spec struct {
		Signals struct {
			Keywords []config.KeywordRule `yaml:"keywords"`
		} `yaml:"signals"`
	} `yaml:"spec"`
}

func testKeywordStandalone(ctx context.Context, client *kubernetes.Clientset, opts testcases.TestCaseOptions) error {
	// Get the absolute path to the CRD file
	// The test is run from the root of the project
	crdPath, err := filepath.Abs("e2e/profiles/dynamic-config/crds/intelligentroute.yaml")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Load the YAML file
	yamlFile, err := ioutil.ReadFile(crdPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal the YAML into our struct
	var route IntelligentRoute
	err = yaml.Unmarshal(yamlFile, &route)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Create a new KeywordClassifier
	classifier, err := classification.NewKeywordClassifier(route.Spec.Signals.Keywords)
	if err != nil {
		return fmt.Errorf("failed to create keyword classifier: %w", err)
	}

	// Test cases
	testCases := []struct {
		Query            string
		ExpectedDecision string
	}{
		{
			Query:            "This is URGENT and needs immediate attention",
			ExpectedDecision: "thinking",
		},
		{
			Query:            "We need this done ASAP",
			ExpectedDecision: "thinking",
		},
		{
			Query:            "urgent: please think about this immediately",
			ExpectedDecision: "thinking",
		},
		{
			Query:            "What is 2 + 2?",
			ExpectedDecision: "",
		},
	}

	var failed bool
	var output strings.Builder

	for _, tc := range testCases {
		decision, _, err := classifier.Classify(tc.Query)
		if err != nil {
			log.Printf("Error classifying query '%s': %v", tc.Query, err)
		}

		if decision != tc.ExpectedDecision {
			failed = true
			output.WriteString(fmt.Sprintf("Query: '%s', Expected: '%s', Actual: '%s'\n", tc.Query, tc.ExpectedDecision, decision))
		}
	}

	if failed {
		return fmt.Errorf("keyword standalone test failed:\n%s", output.String())
	}

	fmt.Println("Keyword standalone test passed!")
	return nil
}
