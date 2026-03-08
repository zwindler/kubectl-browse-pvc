package utils

import (
	"context"
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetNodeTaints(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		client *fake.Clientset
	}{
		{
			name: "Node with Single Taint",
			err:  nil,
			client: fake.NewClientset(
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-node",
					},
					Spec: v1.NodeSpec{
						Taints: []v1.Taint{
							{
								Key:    "key1",
								Value:  "value1",
								Effect: v1.TaintEffectNoSchedule,
							},
						},
					},
				},
			),
		},
		{
			name:   "Node not found returns error",
			err:    nil,
			client: fake.NewClientset(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Node not found returns error" {
				_, err := GetNodeTaints(test.client, "nonexistent-node")
				if err == nil {
					t.Errorf("Expected an error for a nonexistent node, got nil")
				}
				return
			}

			expectedTaints := []v1.Taint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: v1.TaintEffectNoSchedule,
				},
			}

			actualTaints, err := GetNodeTaints(test.client, "test-node")

			if err != nil {
				t.Errorf("Received an err: %s", err)
			} else if !reflect.DeepEqual(actualTaints, expectedTaints) {
				t.Errorf("Expected toleration to be %s, got %s", &expectedTaints[0], &actualTaints[0])
			}

		})
	}

}

func TestBuildTolerations(t *testing.T) {
	tests := []struct {
		name                string
		err                 error
		expectedTolerations []v1.Toleration
		taint               []v1.Taint
	}{
		{
			name: "Single Taint",
			err:  nil,
			taint: []v1.Taint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: v1.TaintEffectNoSchedule,
				},
			},
			expectedTolerations: []v1.Toleration{
				{
					Key:      "key1",
					Value:    "value1",
					Effect:   v1.TaintEffectNoSchedule,
					Operator: v1.TolerationOpEqual,
				},
			},
		},
		{
			name:                "No Taints",
			err:                 nil,
			taint:               []v1.Taint{},
			expectedTolerations: []v1.Toleration{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			actualTolerations := BuildTolerationsForTaints(test.taint)

			if !reflect.DeepEqual(actualTolerations, test.expectedTolerations) {
				t.Errorf("Expected toleration to be %s, got %s", &test.expectedTolerations[0], &actualTolerations[0])
			}
		})
	}
}
