var model = {
  Domains:       [],        // list of domain names
  Catalog:       [],        // list of components
  Components:    [],        // list of component names
  Component:     null,      // the component which is currently being edited
  Controllers:   [],        // list of controllers
  Architectures: {},        //
  Architecture:  null,      // the architecture which is currently being edited
  ArchElement:   null,      // the architectural element which is currently being edited
  Solution:      {},
  SolElement:    null,      // the solution element which is currently being viewed
  Tasks:         [],
  Task:          {},
  Graph:         {          // solution graph
    model         : null,
    view          : null,
    domain        : null ,
    architecture  : null,
    version       : null,
    Components    : {},
    Elements      : {},
    Clusters      : {},
    Relationships : {},
    Nodes         : {},
    Sources       : [],
    Destinations  : [],
    Edges         : {},
    Layers        : [],
    Columns       : [],
    Width         : 0,
    Height        : 0
  },
  Trace:         null,      // task trace
};

//------------------------------------------------------------------------------

// resetModel deletes the current model
function resetModel() {
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/model", {method: "PUT"})
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// saveModel imports a model written in yaml
function saveModel(model) {
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/model", {method: "POST", body: model})
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// loadDomains retrieves a list of domain names from the the repository
function loadDomains() {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/domain")
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Domains = yaml.sort(compareDomains))
}

function compareDomains(a,b) {
  if (a < b) { return -1 }
  if (a > b) { return  1 }
  return 0
}


//------------------------------------------------------------------------------

// saveDomain creates a new domain
function saveDomain(domain) {
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/domain/" + domain, {method: "POST"})
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// deleteDomain removes an existing domain
function deleteDomain(domain) {
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/domain/" + domain, {method: "DELETE"})
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// loadCatalog retrieves a list of components of a domain from the the repository
function loadCatalog(domain) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/catalog/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
    .then((list)     => model.Catalog = list)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a list of component names from the the repository
function loadComponents(domain) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/component/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Components = yaml)
    .then(()         => fetch("http://" + window.location.hostname + ":" + window.location.port + "/controller/" + domain))
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Controllers = yaml)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a component from the the repository
function loadComponent(domain, component) {
  parts   = component.split(" - ")
  name    = parts[0]
  version = parts[1]

  // determine component
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/component/" + domain + "/" + name + "/" + version)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => view.component = yaml)
}

//------------------------------------------------------------------------------

// saveComponent uploads a component to the repository
function saveComponent(domain, comp) {
  body      = jsyaml.safeDump(comp)
  component = comp.Component + " - " + comp.Version
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/component/" + domain, {method: "POST", body: body})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// deleteComponent removes a component from the the repository
function deleteComponent(domain, component) {
  parts   = component.split(" - ")
  name    = parts[0]
  version = parts[1]

  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/component/" + domain + "/" + name + "/" + version, {method: "DELETE"})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// loadArchitectures retrieves a list of architecture names from the the repository
function loadArchitectures(domain) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architectures = yaml.sort())
}

//------------------------------------------------------------------------------

// loadArchitecture retrieves an architecture from the repository
function loadArchitecture(domain, architecture) {
  parts   = architecture.split(" - ")
  name    = parts[0]
  version = parts[1]

  // retrieve architecture
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain + "/" +  name + "/" + version)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architecture = yaml)
}

//------------------------------------------------------------------------------

// saveArchitecture stores an architecture  into the repository
function saveArchitecture(domain, architecture) {
  // save architecture
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain, {
      method: "POST",
      body:   jsyaml.safeDump(architecture)
    })
    .then((response) => response.text())
}

//------------------------------------------------------------------------------

// deployArchitecture deploys or updates a solution based on the architecture
function deployArchitecture(domain, architecture) {
  // deploy architecture
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain + "/" + architecture.Architecture + "/" + architecture.Version, {method: "POST"})
    .then((response) => response.text())
}

//------------------------------------------------------------------------------

// deleteArchitecture removes an architecture  from the repository
function deleteArchitecture(domain, architecture) {
  // delete architecture
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain + "/" + architecture.Architecture + "/" + architecture.Version, {method: "DELETE"})
    .then((response) => response.text())
}

//------------------------------------------------------------------------------

// duplicateArchitecture creates a copy of an architecture with a different name
function duplicateArchitecture(architecture, version) {
  yaml = jsyaml.safeDump(architecture)
  copy = jsyaml.safeLoad(yaml)

  copy.Version = version

  return copy
}

//------------------------------------------------------------------------------

// duplicateElement creates a copy of an element with a different name
function duplicateElement(element, name) {
  yaml = jsyaml.safeDump(element)
  copy = jsyaml.safeLoad(yaml)

  copy.Element = name

  return copy
}

//------------------------------------------------------------------------------

// loadSolutions retrieves a list of solution names from the the repository
function loadSolutions(domain) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/solution/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Solutions = yaml)
}

//------------------------------------------------------------------------------

// loadSolution retrieves a solution from the the repository
function loadSolution(domain, solution) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/solution/" + domain + "/" + solution)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((solution) => {
      // augment clusters and instances with newState, newMin, newMax and newSize
      Object.values(solution.Elements).forEach((e) => {
        Object.values(e.Clusters).forEach((c) => {
          if ( !("NewMin"   in c) ) { c.NewMin   = c.Min   }
          if ( !("NewMax"   in c) ) { c.NewMax   = c.Max   }
          if ( !("NewSize"  in c) ) { c.NewSize  = c.Size  }
          if ( !("NewState" in c) ) { c.NewState = c.State }

          Object.values(c.Instances).forEach((i) => {
            if ( !("NewState" in i) ) { i.NewState = i.State }
          })
        })
      })
      return solution
    })
    .then((solution) => model.Solution = solution)
}

//------------------------------------------------------------------------------

// updateCluster adjusts a cluster regarding dimensions and state
function updateCluster(domain, solution, element, cluster, state, min, max, size) {
  body = jsyaml.safeDump({
    State: state,
    Min:   min,
    Max:   max,
    Size:  size
  })
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/cluster/" + domain + "/" + solution + "/" + element + "/" + cluster, {method: "PUT", body: body})
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// updateInstance adjusts an instance regarding state
function updateInstance(domain, solution, element, cluster, instance, state) {
  body = jsyaml.safeDump({
    State: state
  })
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/instance/" + domain + "/" + solution + "/" + element + "/" + cluster + "/" + instance, {method: "PUT", body: body})
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
}

//------------------------------------------------------------------------------

// loadAll retrieves a solution, the architecture and the catalog from the the repository
function loadAll(domain, solution) {
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/catalog/" + domain)
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
  .then((list)     => model.Catalog = list)
  .then(()         => fetch("http://" + window.location.hostname + ":" + window.location.port + "/solution/" + domain + "/" + solution))
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => model.Solution = yaml)
  .then(()         => fetch("http://" + window.location.hostname + ":" + window.location.port + "/architecture/" + domain + "/" + model.Solution.Solution + " - " +  model.Solution.Version))
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => model.Architecture = yaml)
}

//------------------------------------------------------------------------------

// loadControllers retrieves a list of controllers from the the repository
function loadControllers(domain) {

  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/controller/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Controllers = yaml.sort(compareControllers))
}

function compareControllers(a,b) {
  if (a.Controller < b.Controller) { return -1 }
  if (a.Controller > b.Controller) { return  1 }
  if (a.Version    < b.Version   ) { return -1 }
  if (a.Version    > b.Version   ) { return  1 }
  return 0
}

//------------------------------------------------------------------------------

// addController adds a new controller
function addController(domain, controller) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/controller/" + domain, {
      method: "POST",
      body:   jsyaml.safeDump(controller)
    })
    .then((response) => response.text())
}

//------------------------------------------------------------------------------

// deleteController removes an existing controller
function deleteController(domain, controller) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/controller/" + domain + "/" + controller.Controller + "/" + controller.Version, {
      method: "DELETE",
      body:   jsyaml.safeDump(controller)
    })
    .then((response) => response.text())
}

//------------------------------------------------------------------------------

// loadTasks retrieves a list of task names from the the repository
function loadTasks(domain, solution, element, cluster, instance) {
  // case: only domain and solution are defined
  if (element == "") {
    fetch("http://" + window.location.hostname + ":" + window.location.port + "/tasks/" + domain + "/" + solution)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // case: only domain, solution and element are defined
  if (cluster == "") {
    fetch("http://" + window.location.hostname + ":" + window.location.port + "/tasks/" + domain + "/" + solution + "/" + element)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // case: only domain, solution, element and cluster are defined
  if (instance == "") {
    fetch("http://" + window.location.hostname + ":" + window.location.port + "/tasks/" + domain + "/" + solution + "/" + element + "/" + cluster)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // every parameter is defined
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/tasks/" + domain + "/" + solution + "/" + element + "/" + cluster + "/" + instance)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Tasks = yaml)
}


//------------------------------------------------------------------------------

// loadTrace retrieves a task trace from the the repository
function loadTrace(domain, task) {
  // determine domains
  return fetch("http://" + window.location.hostname + ":" + window.location.port + "/task/" + domain + "/" + task)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Trace = yaml)
}

//------------------------------------------------------------------------------
