package service

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/eriktate/ecs-tasker"
	"github.com/pkg/errors"
)

// A Tasker implements the Tasker interface given an ECS client.
type Tasker struct {
	client *ecs.ECS
}

// UpdateContainerDefinition updates the container definitions for the specified
// task definition. It leaves all other task definition configurations alone.
func (t Tasker) UpdateContainerDefinition(in tasker.UpdateContainerInput) (string, error) {
	taskDef, err := t.getTaskDef(in.Family)
	if err != nil {
		return "", err
	}

	var defs []*ecs.ContainerDefinition
	if err := json.Unmarshal([]byte(in.ContainerDefs), &defs); err != nil { // nolint
		return "", errors.Wrap(err, "could not parse container definitions")
	}

	taskDef.ContainerDefinitions = defs

	rtdr, err := t.client.RegisterTaskDefinition(taskDefToRegisterInput(taskDef))
	if err != nil {
		return "", errors.Wrap(err, "could not register new task definition")
	}

	family := fmt.Sprintf("%s:%d", *rtdr.TaskDefinition.Family, *rtdr.TaskDefinition.Revision)
	return family, nil
}

func (t Tasker) getTaskDef(family string) (*ecs.TaskDefinition, error) {
	req := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(family),
	}

	res, err := t.client.DescribeTaskDefinition(&req)
	if err != nil {
		return nil, errors.Wrap(err, "could not describe task definition")
	}

	return res.TaskDefinition, err
}

func (t Tasker) updateService(cluster, service, family string) error {
	usr := ecs.UpdateServiceInput{
		TaskDefinition: aws.String(family),
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
	}

	if _, err := t.client.UpdateService(&usr); err != nil {
		return errors.Wrap(err, "could not update service")
	}

	return nil
}

func taskDefToRegisterInput(def *ecs.TaskDefinition) *ecs.RegisterTaskDefinitionInput {
	return &ecs.RegisterTaskDefinitionInput{
		Cpu:                     def.Cpu,
		Memory:                  def.Memory,
		ExecutionRoleArn:        def.ExecutionRoleArn,
		Family:                  def.Family,
		ContainerDefinitions:    def.ContainerDefinitions,
		NetworkMode:             def.NetworkMode,
		PlacementConstraints:    def.PlacementConstraints,
		RequiresCompatibilities: def.Compatibilities,
		TaskRoleArn:             def.TaskRoleArn,
		Volumes:                 def.Volumes,
	}
}
