package model

import (
	"fmt"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Cluster
// ========
//
// Attributes:
//   - Version
//   - Target
//   - State
//   - Min
//   - Max
//   - Size
//   - Configuration
//   - Endpoint
//   - Relationships
//   - Instances
//
// Functions:
//   - NewCluster
//
//   - cluster.Show
//   - cluster.Load
//   - cluster.Save
//   - cluster.Update
//   - cluster.Reset
//   - cluster.OK
//
//   - cluster.ListRelationships
//   - cluster.GetRelationship
//   - cluster.AddRelationship
//   - cluster.DeleteRelationship
//
//   - cluster.ListInstances
//   - cluster.GetInstance
//   - cluster.AddInstance
//   - cluster.DeleteInstance
//------------------------------------------------------------------------------

// RelationshipMap is a synchronized map for a map of relationships
type RelationshipMap struct {
	Map          map[string]*Relationship `yaml:"map"` // map of relationships
}

// MarshalYAML marshals a RelationshipMap into yaml
func (m RelationshipMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a RelationshipMap from yaml
func (m *RelationshipMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Relationship{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = RelationshipMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// InstanceMap is a synchronized map for a map of instances
type InstanceMap struct {
	Map          map[string]*Instance     `yaml:"map"` // map of Relationship
}

// MarshalYAML marshals a RelationshipMap into yaml
func (m InstanceMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a RelationshipMap from yaml
func (m *InstanceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Instance{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = InstanceMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Cluster describes the runtime configuration of a solution element cluster within a domain.
type Cluster struct {
	Version       string          `yaml:"version"`       // version of the solution element cluster
	Target        string          `yaml:"target"`        // target state of the solution element cluster
	State         string          `yaml:"state"`         // state of the solution element cluster
	Min           int             `yaml:"min"`           // min. size of the solution element cluster
	Max           int             `yaml:"max"`           // max. size of the solution element cluster
	Size          int             `yaml:"size"`          // size of the solution element cluster
	Configuration string          `yaml:"configuration"` // runtime configuration of the solution element cluster
	Endpoint      string          `yaml:"endpoint"`      // endpoint of the solution element cluster
	Relationships RelationshipMap `yaml:"relationships"` // relationships of the solution element cluster
	Instances     InstanceMap     `yaml:"instances"`     // instances of the solution element cluster
}

//------------------------------------------------------------------------------

// NewCluster creates a new cluster
func NewCluster(version string, state string, min int, max int, size int, configuration string) (*Cluster, error) {
	var cluster Cluster

	cluster.Version       = version
	cluster.Target        = state
	cluster.State         = InitialState
	cluster.Min           = min
	cluster.Max           = max
	cluster.Size          = size
	cluster.Configuration = configuration
	cluster.Endpoint      = ""
	cluster.Relationships = RelationshipMap{Map: map[string]*Relationship{}}
	cluster.Instances     = InstanceMap{Map: map[string]*Instance{}}

	// success
	return &cluster, nil
}

//------------------------------------------------------------------------------

// Show displays the cluster information as yaml
func (cluster *Cluster) Show() (string, error) {
	return util.ConvertToYAML(cluster)
}

//------------------------------------------------------------------------------

// Save writes the cluster as yaml data to a file
func (cluster *Cluster) Save(filename string) error {
	return util.SaveYAML(filename, cluster)
}

//------------------------------------------------------------------------------

// Load reads the cluster from a file
func (cluster *Cluster) Load(filename string) error {
	return util.LoadYAML(filename, cluster)
}

//------------------------------------------------------------------------------

// ListRelationships lists the names of all defined relationships
func (cluster *Cluster) ListRelationships() ([]string, error) {
	// collect names
	relationships := []string{}

	for relationship := range cluster.Relationships.Map {
		relationships = append(relationships, relationship)
	}

	// success
	return relationships, nil
}

//------------------------------------------------------------------------------

// GetRelationship retrieves a relationship by name
func (cluster *Cluster) GetRelationship(name string) (*Relationship, error) {
	// determine relationship
	relationship, ok := cluster.Relationships.Map[name]

	if !ok {
		return nil, errors.New("relationship not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------

// AddRelationship adds a relationship to a cluster
func (cluster *Cluster) AddRelationship(relationship *Relationship) {
	cluster.Relationships.Map[relationship.Relationship] = relationship
}

//------------------------------------------------------------------------------

// DeleteRelationship deletes a relationship from a cluster
func (cluster *Cluster) DeleteRelationship(name string) {
	delete(cluster.Relationships.Map, name)
}

//------------------------------------------------------------------------------

// ListInstances lists the names of all defined instances
func (cluster *Cluster) ListInstances() ([]string, error) {
	// collect names
	instances := []string{}

	for instance := range cluster.Instances.Map {
		instances = append(instances, instance)
	}

	// success
	return instances, nil
}

//------------------------------------------------------------------------------

// GetInstance retrieves an instance by uuid
func (cluster *Cluster) GetInstance(uuid string) (*Instance, error) {
	// determine instance
	instance, ok := cluster.Instances.Map[uuid]

	if !ok {
		return nil, errors.New("instance not found")
	}

	// success
	return instance, nil
}

//------------------------------------------------------------------------------

// AddInstance adds an instance to a cluster
func (cluster *Cluster) AddInstance(instance *Instance) {
	cluster.Instances.Map[instance.UUID] = instance
}

//------------------------------------------------------------------------------

// DeleteInstance deletes an instance from a cluster
func (cluster *Cluster) DeleteInstance(uuid string) {
	delete(cluster.Instances.Map, uuid)
}

//------------------------------------------------------------------------------

// Update instantiates/update a cluster based on a cluster configuration.
func (cluster *Cluster) Update(domainName string, solutionName string, elementName string, clusterConfiguration *ClusterConfiguration) error {
	// check if the names are compatible
	if cluster.Version != clusterConfiguration.Version {
		return errors.New("Version of cluster does match the version defined in the cluster configuration")
	}

	// update target state and sizes
	cluster.Target = clusterConfiguration.State
	cluster.Min    = clusterConfiguration.Min
	cluster.Max    = clusterConfiguration.Max
	cluster.Size   = clusterConfiguration.Size

	// check compatability of all relationships
	relationshipNames, _ := clusterConfiguration.ListRelationships()
	for _, relationshipName := range relationshipNames {

		relationship, _              := cluster.GetRelationship(relationshipName)
		relationshipConfiguration, _ := clusterConfiguration.GetRelationship(relationshipName)

		// relationship already exists
		if cluster != nil {
			// check compatability of references
			if relationship.Element != relationshipConfiguration.Element ||
			   relationship.Version != relationshipConfiguration.Version {
					 return fmt.Errorf("Incompatible relationship: '%s' of the cluster '%s'", relationshipName, cluster.Version)
				 }
		} else {
			// relationship does not exist
			// create new relationship
			relationship, _ = NewRelationship(relationshipName,
				                                relationshipConfiguration.Dependency,
																				relationshipConfiguration.Type,
																				domainName,
																				solutionName,
																				relationshipConfiguration.Element,
																				relationshipConfiguration.Version,
																				"")
			cluster.AddRelationship(relationship)
		}
	}

	// add missing instances
	instanceNames, _ := cluster.ListInstances()
	targetSize       := clusterConfiguration.Max
	currentSize      := len(instanceNames)
	for currentSize < targetSize {
		// add new instance to cluster in its initial state
		instance, _ := NewInstance(util.UUID(), InitialState, "")
		cluster.AddInstance(instance)

		currentSize = currentSize + 1
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Reset state of cluster
func (cluster *Cluster) Reset() {
	cluster.Target = InitialState

	// reset all relationships
	relationshipNames, _ := cluster.ListRelationships()
	for _, relationshipName := range relationshipNames {
		relationship, _ := cluster.GetRelationship(relationshipName)

		relationship.Reset()
	}

	// reset all instances
	instanceNames, _ := cluster.ListInstances()
	for _, instanceName := range instanceNames {
		instance, _ := cluster.GetInstance(instanceName)

		instance.Reset()
	}
}

//------------------------------------------------------------------------------

// OK checks if the cluster has converged to the target state
func (cluster *Cluster) OK() bool {
	// check state
	if cluster.Target != cluster.State {
		return false
	}

	// check size
	instanceNames, _ := cluster.ListInstances()
	instances := len(instanceNames)
	if instances != cluster.Size {
		return false
	}

	// check relationships
	switch cluster.State {
	case InactiveState:
		// check if all context relationships are active
		relationshipNames, _ := cluster.ListRelationships()
		for _, relationshipName := range relationshipNames {
			relationship, _ := cluster.GetRelationship(relationshipName)

			if relationship.Type != ContextRelationship {
				continue
			}

			refCluster, _ := GetCluster(relationship.Domain, relationship.Solution, relationship.Element, relationship.Version)
			if refCluster.State != ActiveState	{
				return false
			}
		}
	case ActiveState:
		// check if all service and context relationships are active
		relationshipNames, _ := cluster.ListRelationships()
		for _, relationshipName := range relationshipNames {
			relationship, _ := cluster.GetRelationship(relationshipName)

			if relationship.Type != ContextRelationship && relationship.Type != ServiceRelationship {
				continue
			}

			refCluster, _ := GetCluster(relationship.Domain, relationship.Solution, relationship.Element, relationship.Version)
			if refCluster.State != ActiveState	{
				return false
			}
		}
	}

	// cluster is ok
	return true
}

//------------------------------------------------------------------------------
