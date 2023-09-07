// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reactor

import (
	"testing"

	"github.com/gardener/docforge/pkg/manifestadapter"
	"github.com/gardener/docforge/pkg/resourcehandlers/resourcehandlersfakes"

	"github.com/gardener/docforge/pkg/resourcehandlers"
	"github.com/gardener/docforge/pkg/util/tests"
)

func init() {
	tests.SetKlogV(6)
}

var (
	apiRefNode = &manifestadapter.Node{
		Name: "apiRef",
	}

	archNode = &manifestadapter.Node{
		Name: "arch",
		Nodes: []*manifestadapter.Node{
			apiRefNode,
		},
	}

	blogNode = &manifestadapter.Node{
		Name: "blog",
	}

	tasksNode = &manifestadapter.Node{
		Name: "tasks",
	}
)

func createNewDocumentation() *manifestadapter.Documentation {
	return &manifestadapter.Documentation{
		Structure: []*manifestadapter.Node{
			{
				Name: "rootNode",
				Nodes: []*manifestadapter.Node{
					archNode,
					blogNode,
					tasksNode,
				},
			},
		},
	}
}

func Test_tasks(t *testing.T) {
	newDoc := createNewDocumentation()
	type args struct {
		node  *manifestadapter.Node
		tasks []interface{}
		// lds   localityDomain
	}
	tests := []struct {
		name          string
		args          args
		expectedTasks []*DocumentWorkTask
	}{
		{
			name: "it creates tasks based on the provided doc",
			args: args{
				node:  newDoc.Structure[0],
				tasks: []interface{}{},
			},
			expectedTasks: []*DocumentWorkTask{
				{
					Node: newDoc.Structure[0],
				},
				{
					Node: archNode,
				},
				{
					Node: apiRefNode,
				},
				{
					Node: blogNode,
				},
				{
					Node: tasksNode,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeRH := resourcehandlersfakes.FakeResourceHandler{}
			rhs := resourcehandlers.NewRegistry(&fakeRH)
			tasks([]*manifestadapter.Node{tc.args.node}, &tc.args.tasks)

			if len(tc.args.tasks) != len(tc.expectedTasks) {
				t.Errorf("expected number of tasks %d != %d", len(tc.expectedTasks), len(tc.args.tasks))
			}

			for i, task := range tc.args.tasks {
				if task.(*DocumentWorkTask).Node.Name != tc.expectedTasks[i].Node.Name {
					t.Errorf("expected task with Node name %s != %s", task.(*DocumentWorkTask).Node.Name, tc.expectedTasks[i].Node.Name)
				}
			}
			rhs.Remove()
		})
	}
}
