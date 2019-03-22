package model

import (
	"fmt"
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
//   - cluster.Pools
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

// Cluster describes the runtime configuration of a solution element cluster within a domain.
type Cluster struct {
	Version        string                   `yaml:"Version"`                  // version of the solution element cluster
	Target         string                   `yaml:"Target"`                   // target state of the solution element cluster
	State          string                   `yaml:"State"`                    // state of the solution element cluster
	Min            int                      `yaml:"Min"`                      // min. size of the solution element cluster
	Max            int                      `yaml:"Max"`                      // max. size of the solution element cluster
	Size           int                      `yaml:"Size"`                     // size of the solution element cluster
	Configuration  string                   `yaml:"Configuration"`            // runtime configuration of the solution element cluster
	Endpoint       string                   `yaml:"Endpoint"`                 // endpoint of the solution element cluster
	Relationships  map[string]*Relationship `yaml:"Relationships"`            // relationships of the solution element cluster
	RelationshipsX sync.RWMutex             `yaml:"RelationshipsX,omitempty"` // mutex for relationships
	Instances      map[string]*Instance     `yaml:"Instances"`                // instances of the solution element cluster
	InstancesX     sync.RWMutex             `yaml:"InstancesX,omitempty"`     // mutex for instances
}

//------------------------------------------------------------------------------

// NewCluster creates a new cluster
func NewCluster(version string, state string, min int, max int, size int, configuration string) (*Cluster, error) {
	var cluster Cluster

	cluster.Version        = version
	cluster.Target         = state
	cluster.State          = InitialState
	cluster.Min            = min
	cluster.Max            = max
	cluster.Size           = size
	cluster.Configuration  = configuration
	cluster.Endpoint       = ""
	cluster.Relationships  = map[string]*Relationship{}
	cluster.RelationshipsX = sync.RWMutex{}
	cluster.Instances      = map[string]*Instance{}
	cluster.InstancesX     = sync.RWMutex{}

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

  cluster.RelationshipsX.RLock()
	for relationship := range cluster.Relationships {
		relationships = append(relationships, relationship)
	}
	cluster.RelationshipsX.RUnlock()

	// success
	return relationships, nil
}

//------------------------------------------------------------------------------

// GetRelationship retrieves a relationship by name
func (cluster *Cluster) GetRelationship(name string) (*Relationship, error) {
	// determine relationship
	cluster.RelationshipsX.RLock()
	relationship, ok := cluster.Relationships[name]
	cluster.RelationshipsX.RUnlock()

	if !ok {
		return nil, errors.New("relationship not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------

// AddRelationship adds a relationship to a cluster
func (cluster *Cluster) AddRelationship(relationship *Relationship) {
	cluster.RelationshipsX.Lock()
	cluster.Relationships[relationship.Relationship] = relationship
	cluster.RelationshipsX.Unlock()
}

//------------------------------------------------------------------------------

// DeleteRelationship deletes a relationship from a cluster
func (cluster *Cluster) DeleteRelationship(name string) {
	cluster.RelationshipsX.Lock()
	delete(cluster.Relationships, name)
	cluster.RelationshipsX.Unlock()
}

//------------------------------------------------------------------------------

// ListInstances lists the names of all defined instances
func (cluster *Cluster) ListInstances() ([]string, error) {
	// collect names
	instances := []string{}

  cluster.InstancesX.RLock()
	for instance := range cluster.Instances {
		instances = append(instances, instance)
	}
	cluster.InstancesX.RUnlock()

	// success
	return instances, nil
}

//------------------------------------------------------------------------------

// GetInstance retrieves an instance by uuid
func (cluster *Cluster) GetInstance(uuid string) (*Instance, error) {
	// determine instance
	cluster.InstancesX.RLock()
	instance, ok := cluster.Instances[uuid]
	cluster.InstancesX.RUnlock()

	if !ok {
		return nil, errors.New("instance not found")
	}

	// success
	return instance, nil
}

//------------------------------------------------------------------------------

// AddInstance adds an instance to a cluster
func (cluster *Cluster) AddInstance(instance *Instance) {
	cluster.InstancesX.Lock()
	cluster.Instances[instance.UUID] = instance
	cluster.InstancesX.Unlock()
}

//------------------------------------------------------------------------------

// DeleteInstance deletes an instance from a cluster
func (cluster *Cluster) DeleteInstance(uuid string) {
	cluster.InstancesX.Lock()
	delete(cluster.Instances, uuid)
	cluster.InstancesX.Unlock()
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
		if relationship != nil {
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

	// check size of cluster
	// count number of instances in each lifecycle state
	_, inactive, active, failure, _ := cluster.Pools()

	// check if there are any failed instances
	if failure > 0 {
		return false
	}

	// check categories
	switch cluster.Target {
		case InitialState:
			// too many instances are still active or have been deployed
			if inactive > 0 || active > 0 {
				return false
			}
		case InactiveState:
			// too many instances are still active
			if active > 0 {
				return false
			}
			// size of cluster does not match
			if inactive != cluster.Size {
				return false
			}
		case ActiveState:
			// size of cluster does not match
			if active != cluster.Size {
				return false
			}
			// size of inactive instances is still too high
			if inactive > (cluster.Max - cluster.Min) {
				return false
			}
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

// Pools determines the sizes of the pools in which the instances may reside
func (cluster *Cluster) Pools() (initial int, inactive int, active int, failure int, other int) {
	// count number of instances in each lifecycle state
	initial  = 0
	inactive = 0
	active   = 0
	failure  = 0
	other    = 0

	instanceNames, _ := cluster.ListInstances()
	for _, instanceName := range instanceNames {
		instance, _ := cluster.GetInstance(instanceName)
		switch instance.State {
			case InitialState:
				initial++
			case InactiveState:
				inactive++
			case ActiveState:
				active++
			case FailureState:
				failure++
			default:
				other++
		}
	}

	// finished
	return initial, inactive, active, failure, other
}

//------------------------------------------------------------------------------
