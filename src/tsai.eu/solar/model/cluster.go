package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Cluster
// ========
//
// Attributes:
//   - Version
//   - State
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
	sync.RWMutex `yaml:"mutex,omitempty"`              // mutex
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
	sync.RWMutex `yaml:"mutex,omitempty"`              // mutex
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
	State         string          `yaml:"state"`         // state of the solution element cluster
	Size          int             `yaml:"size"`          // size of the solution element cluster
	Configuration string          `yaml:"configuration"` // runtime configuration of the solution element cluster
	Endpoint      string          `yaml:"endpoint"`      // endpoint of the solution element cluster
	Relationships RelationshipMap `yaml:"relationships"` // relationships of the solution element cluster
	Instances     InstanceMap     `yaml:"instances"`     // instances of the solution element cluster
}

//------------------------------------------------------------------------------

// NewCluster creates a new cluster
func NewCluster(version string, state string, size int, configuration string) (*Cluster, error) {
	var cluster Cluster

	cluster.Version = version
	cluster.State = state
	cluster.Size = size
	cluster.Configuration = configuration
	cluster.Endpoint = ""
	cluster.Relationships = RelationshipMap{Map: map[string]*Relationship{}}
	cluster.Instances = InstanceMap{Map: map[string]*Instance{}}

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

	cluster.Relationships.RLock()
	for relationship := range cluster.Relationships.Map {
		relationships = append(relationships, relationship)
	}
	cluster.Relationships.RUnlock()

	// success
	return relationships, nil
}

//------------------------------------------------------------------------------

// GetRelationship retrieves a relationship by name
func (cluster *Cluster) GetRelationship(name string) (*Relationship, error) {
	// determine relationship
	cluster.Relationships.RLock()
	relationship, ok := cluster.Relationships.Map[name]
	cluster.Relationships.RUnlock()

	if !ok {
		return nil, errors.New("relationship not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------

// AddRelationship adds a relationship to a cluster
func (cluster *Cluster) AddRelationship(relationship *Relationship) {
	cluster.Relationships.Lock()
	cluster.Relationships.Map[relationship.Relationship] = relationship
	cluster.Relationships.Unlock()
}

//------------------------------------------------------------------------------

// DeleteRelationship deletes a relationship from a cluster
func (cluster *Cluster) DeleteRelationship(name string) {
	cluster.Relationships.Lock()
	delete(cluster.Relationships.Map, name)
	cluster.Relationships.Unlock()
}

//------------------------------------------------------------------------------

// ListInstances lists the names of all defined instances
func (cluster *Cluster) ListInstances() ([]string, error) {
	// collect names
	instances := []string{}

	cluster.Instances.RLock()
	for instance := range cluster.Instances.Map {
		instances = append(instances, instance)
	}
	cluster.Instances.RUnlock()

	// success
	return instances, nil
}

//------------------------------------------------------------------------------

// GetInstance retrieves an instance by uuid
func (cluster *Cluster) GetInstance(uuid string) (*Instance, error) {
	// determine instance
	cluster.Instances.RLock()
	instance, ok := cluster.Instances.Map[uuid]
	cluster.Instances.RUnlock()

	if !ok {
		return nil, errors.New("instance not found")
	}

	// success
	return instance, nil
}

//------------------------------------------------------------------------------

// AddInstance adds an instance to a cluster
func (cluster *Cluster) AddInstance(instance *Instance) {
	cluster.Instances.Lock()
	cluster.Instances.Map[instance.UUID] = instance
	cluster.Instances.Unlock()
}

//------------------------------------------------------------------------------

// DeleteInstance deletes an instance from a cluster
func (cluster *Cluster) DeleteInstance(uuid string) {
	cluster.Instances.Lock()
	delete(cluster.Instances.Map, uuid)
	cluster.Instances.Unlock()
}

//------------------------------------------------------------------------------
