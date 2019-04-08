package model

import (
	"errors"
	"strings"
	"strconv"
	"regexp"

	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Relationship
// ============
//
// Attributes:
//   - Relationship
//   - Dependency
//   - Type
//   - Element
//   - Version
//   - Target
//   - State
//   - Configuration
//   - Endpoint
//
// Functions:
//   - NewRelationship
//
//   - relationship.Show
//   - relationship.Load
//   - relationship.Save
//   - relationship.Reset
//------------------------------------------------------------------------------

// Relationship describes the runtime configuration of a relationship between clusters within a domain.
type Relationship struct {
	Relationship  string  `yaml:"Relationship"`  // name of the relationship
	Dependency    string  `yaml:"Dependency"`    // name of the dependency
	Type          string  `yaml:"Type"`          // type of dependency
	Domain        string  `yaml:"Domain"`        // domain to which this relationship refers to
	Solution      string  `yaml:"Solution"`      // solution to which this relationship refers to
	Element       string  `yaml:"Element"`       // element to which this relationship refers to
	Version       string  `yaml:"Version"`       // version of the element to which this relationship refers to
	Target        string  `yaml:"Target"`        // target state of relationship
	State         string  `yaml:"State"`         // current state of relationship
	Configuration string  `yaml:"Configuration"` // runtime configuration of the relationship
	Endpoint      string  `yaml:"Endpoint"`      // endpoint of the relationship
}

//------------------------------------------------------------------------------

// NewRelationship creates a new relationship
func NewRelationship(name string, dependency string, dependencyType string, domain string, solution string, element string, version string, configuration string) (*Relationship, error) {
	var relationship Relationship

	relationship.Relationship  = name
	relationship.Dependency    = dependency
	relationship.Type          = dependencyType
	relationship.Domain        = domain
	relationship.Solution      = solution
	relationship.Element       = element
	relationship.Version       = version
	relationship.Target        = InitialState
	relationship.State         = InitialState
	relationship.Configuration = configuration
	relationship.Endpoint      = ""

	// success
	return &relationship, nil
}

//------------------------------------------------------------------------------

// renderConfiguration calculates the configuration from the component template and the parameters defined in the relationshipConfiguration.
func (relationship *Relationship) renderConfiguration(domainName string, solutionName string, version string, element *Element, cluster *Cluster, relationshipConfiguration *RelationshipConfiguration) {
	// determine component
	component, err := GetComponent(domainName, element.Component + " - " + cluster.Version)
	if err != nil {
		util.LogError("relationship", "MODEL", "unknown component '" + element.Component + " - " + cluster.Version + "' within domain: '" + domainName + "'")
		return
	}

	// determine dependency
	dependency, err := component.GetDependency(relationshipConfiguration.Dependency)
	if err != nil {
		util.LogError("relationship", "MODEL", "unknown dependency '" + element.Component + " - " + cluster.Version + " / " + relationshipConfiguration.Relationship + "' within domain: '" + domainName + "'")
		return
	}

	// get parameters
	parameters := map[string]string{}
	err = util.ConvertFromYAML(relationshipConfiguration.Configuration, &parameters)
	if err != nil {
		util.LogError("relationship", "MODEL", "unable to parse the parameters defined in the architecture cluster relationship: '" + element.Element + " - " + cluster.Version + " / " + relationshipConfiguration.Relationship + "' within domain: '" + domainName + "'")
	}
	if len(parameters) == 0 {
		parameters = map[string]string{}
	}

	// add default parameters
	parameters["domain"]       = domainName
	parameters["solution"]     = solutionName
	parameters["version"]      = version
	parameters["element"]      = element.Element
	parameters["component"]    = element.Component
	parameters["cluster"]      = cluster.Version
	parameters["min"]          = strconv.Itoa(cluster.Min)
	parameters["max"]          = strconv.Itoa(cluster.Max)
	parameters["size"]         = strconv.Itoa(cluster.Size)
	parameters["relationship"] = relationshipConfiguration.Relationship

	// determine all required parameters
	configuration := dependency.Configuration
	r := regexp.MustCompile(`{{([^}]*)}}`)
	matches := r.FindAllStringSubmatch(configuration, -1)
	if matches != nil {
		for _, match := range matches {
			name := match[1]
			key  := strings.TrimSpace(name)

			value, ok := parameters[key]
			if ok {
				configuration = strings.Replace(configuration, "{{" + name + "}}", value, -1)
			}
		}
	}

	// set conifguration of relationship
	relationship.Configuration = configuration
}

//------------------------------------------------------------------------------

// Update instantiates/update a relationship based on a relationship configuration.
func (relationship *Relationship) Update(domainName string, solutionName string, version string, element *Element, cluster *Cluster, relationshipConfiguration *RelationshipConfiguration) error {
	// check if the names are compatible
	if relationship.Relationship != relationshipConfiguration.Relationship {
		return errors.New("Name of relationship does not match the name defined in the relationship configuration")
	}

	// update configuration
	relationship.renderConfiguration(domainName, solutionName, version, element, cluster, relationshipConfiguration)

	// success
	return nil
}

//------------------------------------------------------------------------------

// Show displays the relationship information as yaml
func (relationship *Relationship) Show() (string, error) {
	return util.ConvertToYAML(relationship)
}

//------------------------------------------------------------------------------

// Save writes the relationship as yaml data to a file
func (relationship *Relationship) Save(filename string) error {
	return util.SaveYAML(filename, relationship)
}

//------------------------------------------------------------------------------

// Load reads the relationship from a file
func (relationship *Relationship) Load(filename string) error {
	return util.LoadYAML(filename, relationship)
}

//------------------------------------------------------------------------------

// Reset state of relationship
func (relationship *Relationship) Reset() {
	relationship.Target = InitialState
}

//------------------------------------------------------------------------------
