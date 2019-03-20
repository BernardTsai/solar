var model = {
  Domains:       [],        // list of domain names
  Catalog:       [],        // list of components
  Components:    [],        // list of component names
  Component:     null,      // the component which is currently being edited
  Architectures: {},        //
  Architecture:  null,      // the architecture which is currently being edited
  ArchElement:   null,      // the architectural element which is currently being edited
  Solution:      {},
  SolElement:    null,      // the solution element which is currently being viewed
  Tasks:         [],
  Task:          {},
  Graph:         null,      // solution graph
  Trace:         null,      // task trace
};

//------------------------------------------------------------------------------

// loadDomains retrieves a list of domain names from the the repository
function loadDomains() {
  // determine domains
  return fetch("http://localhost/domain")
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Domains = yaml)
}


//------------------------------------------------------------------------------

// loadCatalog retrieves a list of components of a domain from the the repository
function loadCatalog(domain) {
  // determine domains
  return fetch("http://localhost/catalog/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
    .then((list)     => model.Catalog = list)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a list of component names from the the repository
function loadComponents(domain) {
  // determine domains
  return fetch("http://localhost/component/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Components = yaml)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a component from the the repository
function loadComponent(domain, component) {
  // determine component
  return fetch("http://localhost/component/" + domain + "/" + component)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => view.component = yaml)
}

//------------------------------------------------------------------------------

// saveComponent uploads a component to the repository
function saveComponent(domain, comp) {
  body      = jsyaml.safeDump(comp)
  component = comp.Component + " - " + comp.Version
  return fetch("http://localhost/component/" + domain, {method: "POST", body: body})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// deleteComponent removes a component from the the repository
function deleteComponent(domain, component) {
  return fetch("http://localhost/component/" + domain + "/" + component, {method: "DELETE"})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// loadArchitectures retrieves a list of architecture names from the the repository
function loadArchitectures(domain) {
  // determine domains
  return fetch("http://localhost/architecture/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architectures = yaml)
}

//------------------------------------------------------------------------------

// loadArchitecture retrieves an  architecture from the the repository
function loadArchitecture(domain, architecture) {
  // determine domains
  return fetch("http://localhost/architecture/" + domain + "/" + architecture)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architecture = yaml)
}

//------------------------------------------------------------------------------

// loadSolutions retrieves a list of solution names from the the repository
function loadSolutions(domain) {
  // determine domains
  return fetch("http://localhost/solution/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Solutions = yaml)
}

//------------------------------------------------------------------------------

// loadSolution retrieves a solution from the the repository
function loadSolution(domain, solution) {
  // determine domains
  return fetch("http://localhost/solution/" + domain + "/" + solution)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Solution = yaml)
}

//------------------------------------------------------------------------------

// loadAll retrieves a solution, the architecture and the catalog from the the repository
function loadAll(domain, solution) {
  // determine domains
  return fetch("http://localhost/catalog/" + domain)
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
  .then((list)     => model.Catalog = list)
  .then(()         => fetch("http://localhost/solution/" + domain + "/" + solution))
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => model.Solution = yaml)
  .then(()         => fetch("http://localhost/architecture/" + domain + "/" + model.Solution.Solution + " - " +  model.Solution.Version))
  .then((response) => response.text())
  .then((text)     => jsyaml.safeLoad(text))
  .then((yaml)     => model.Architecture = yaml)
}

//------------------------------------------------------------------------------

// loadTasks retrieves a list of task names from the the repository
function loadTasks(domain, solution, element, cluster, instance) {
  // case: only domain and solution are defined
  if (element == "") {
    fetch("http://localhost/tasks/" + domain + "/" + solution)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // case: only domain, solution and element are defined
  if (cluster == "") {
    fetch("http://localhost/tasks/" + domain + "/" + solution + "/" + element)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // case: only domain, solution, element and cluster are defined
  if (instance == "") {
    fetch("http://localhost/tasks/" + domain + "/" + solution + "/" + element + "/" + cluster)
      .then((response) => response.text())
      .then((text)     => jsyaml.safeLoad(text))
      .then((yaml)     => model.Tasks = yaml)
    return
  }

  // every parameter is defined
  return fetch("http://localhost/tasks/" + domain + "/" + solution + "/" + element + "/" + cluster + "/" + instance)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Tasks = yaml)
}


//------------------------------------------------------------------------------

// loadTrace retrieves a task trace from the the repository
function loadTrace(domain, task) {
  // determine domains
  return fetch("http://localhost/task/" + domain + "/" + task)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Trace = yaml)
}

//------------------------------------------------------------------------------
