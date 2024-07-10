package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	operatorExporterPackage "github.com/krateoplatformops/finops-operator-exporter/api/v1"
	operatorPackage "github.com/krateoplatformops/finops-operator-vm-manager/api/v1"
	providers "github.com/krateoplatformops/finops-operator-vm-manager/providers"
	prometheusExporterGeneric "github.com/krateoplatformops/finops-prometheus-exporter-generic/pkg/config"
)

type OptimizationRequest struct {
	ResourceId   string       `json:"resourceId"`
	Optimization Optimization `json:"optimization"`
}

type Optimization struct {
	ResourceName  string     `json:"resourceName"`
	ResourceDelta int        `json:"resourceDelta"`
	TypeChange    TypeChange `json:"typeChange"`
}

type TypeChange struct {
	Cyclic string `json:"cyclic"`
	From   string `json:"from"`
	To     string `json:"to"`
}

func Fatal(err error) {
	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
	}
}

func GetClientSet() (*kubernetes.Clientset, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return &kubernetes.Clientset{}, err
	}

	inClusterConfig.APIPath = "/apis"
	inClusterConfig.GroupVersion = &operatorPackage.GroupVersion

	clientset, err := kubernetes.NewForConfig(inClusterConfig)
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	return clientset, nil
}

func CreateOptimizationCustomResource(p *OptimizationRequest, optName string, optNamespace string, secretName string, secretNamespace string) error {
	clientset, err := GetClientSet()
	if err != nil {
		return err
	}

	action := "nop"
	if p.Optimization.ResourceDelta < 0 {
		action = "scale-down"
	} else if p.Optimization.ResourceDelta > 0 {
		action = "scale-up"
	}

	// Check if the ConfigManagerVM already exists
	jsonData, _ := clientset.RESTClient().Get().
		AbsPath("/apis/finops.krateo.io/v1").
		Namespace(optNamespace).
		Resource("configmanagervms").
		Name(optName).
		DoRaw(context.TODO())

	var crdResponse prometheusExporterGeneric.Kind
	_ = json.Unmarshal(jsonData, &crdResponse)
	// If it does not exist, build it
	if crdResponse.Kind == "Status" && crdResponse.Status == "Failure" {
		configManagerVM := operatorPackage.ConfigManagerVM{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigManagerVM",
				APIVersion: "finops.krateo.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      optName,
				Namespace: optNamespace,
			},
			Spec: operatorPackage.ConfigManagerVMSpec{
				ResourceProvider: "azure",
				ProviderSpecificResources: operatorPackage.ProviderSpecificResources{
					AzureLogin: providers.Azure{
						TokenRef: operatorExporterPackage.ObjectRef{
							Name:      secretName,
							Namespace: secretNamespace,
						},
						Path:          "/subscriptions/31a6f7f0-996a-490e-b431-25c3f54b08b5/resourcegroups/FinOps/providers/Microsoft.Compute/virtualMachines/poppy",
						ResourceDelta: p.Optimization.ResourceDelta,
						Action:        action,
					},
				},
			},
		}
		jsonData, err = json.Marshal(configManagerVM)
		if err != nil {
			return err
		}
		// Create the object in the cluster
		_, err := clientset.RESTClient().Post().
			AbsPath("/apis/finops.krateo.io/v1").
			Namespace(optNamespace).
			Resource("configmanagervms").
			Name(optName).
			Body(jsonData).
			DoRaw(context.TODO())

		if err != nil {
			return err
		}
	}
	return nil
}
