var view = {
  nav:          "Components",    // selected view (Components,Architecture,Solution,Automation)
  domain:       "",              // selected domain
  component:    -1,              // index of selected component
  newComponent: false,           // indicates if a new component is to be designed
  solution:     "",              // name of architecture/solution
  version:      "",              // version of architecture/solution
  cluster:      "",              // name of selected cluster
  relationship: "",              // name of selected relationship

  ce: {                          // fields for the component editor
    New:            false,
    Component:      "",
    Version:        "",
    Configuration1: "",
    Dependency:     "",
    DepType:        "",
    DepComponent:   "",
    DepVersion:     "",
    Configuration2: ""
  },

  ae: {                          // fields for the architecture element editor
    New:            false,
    Element:        "unknown",
    Component:      "",
    Configuration1: "",
    Cluster:        "",
    State:          "initial",
    Min:            "1",
    Max:            "1",
    Size:           "1",
    Configuration2: "",
    Relationship:   "",
    Dependency:     "",
    DepType:        "",
    RelElement:     "",
    Configuration3: ""
  },

  se: {                          // fields for the solution element editor
    New:            false,
    Element:        "unknown",
    Component:      "",
    Configuration1: "",
    Cluster:        "",
    Target:         "initial",
    State:          "initial",
    Min:            "1",
    Max:            "1",
    Size:           "1",
    Configuration2: "",
    Relationship:   "",
    Dependency:     "",
    DepType:        "",
    RelElement:     "",
    Configuration3: ""
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
    }
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
