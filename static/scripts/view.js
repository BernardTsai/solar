var view = {
  nav:          "Administration",  // selected view (Components,Architecture,Solution,Automation,Administration)
  subnav:       "Model",           // selected subview (Model,Logs, Controller for Administration)
  modelDomain:  "",                // selected domain within the administration view
  domain:       "",                // selected domain
  component:    -1,                // index of selected component
  newComponent: false,             // indicates if a new component is to be designed
  architecture: "",                // name of architecture
  solution:     "",                // name of architecture/solution
  version:      "",                // version of architecture/solution
  element:      "",                // name of selected element
  cluster:      "",                // name of selected cluster
  relationship: "",                // name of selected relationship
  showGraph:    true,              // specifies wether to show graph or editor

  ce: {                          // fields for the component editor
    New:            false,
    Component:      "",
    Version:        "",
    Dependency:     "",
    Configuration:  null,
    ConfTitle:      ""
  },

  ae: {                    // fields for the architecture element editor
    ConfType:       "",    // type of configuration: element, cluster, relationship
    ConfTitle:      "",    // title of configuration editor
    Display:        "",    // display paramters, template or configuration in configuration editor
    Configuration:  null,  // holds the configuration parameters information
    Cluster:        "",    // the name of the cluster information in focus
    Relationship:   "",    // the name of the relationship information in focus
  },

  se: {                    // fields for the solution element editor
    ConfType:       "",    // type of configuration: element, cluster, relationship
    ConfTitle:      "",    // title of configuration editor
    Configuration:  null,  // holds the configuration parameters information
    Cluster:        "",    // the name of the cluster information in focus
    Relationship:   "",    // the name of the relationship information in focus
  },

  graph: {
    dx: 40,
    dy: 40,
    node: {
      width:  160,
      height:  40
    },
    port: {
      diameter: 8,
      border:   1
    },
    viewElement: null
  },

  automation: {
    solution: "",
    element:  "",
    cluster:  "",
    instance: "",
    task:     "",
    line:     16,
    width:    1000
  },

  instance:  "",
  focus:     "",
  node:      null
}
