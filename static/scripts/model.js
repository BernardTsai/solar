var model = {
  Domains:       [],        // list of domain names
  Catalog:       [],        // list of components
  Components:    [],        // list of component names
  Component:     null,      // the component which is currently being edited
  Architectures: {},        //
  Architecture:  null,      // the architecture which is currently being edited
  ArchElement:   null,      // the architectural element which is currently being edited
  Solution:      {},
  Tasks:         [],
  Task:          {}
};

//------------------------------------------------------------------------------

// loadDomains retrieves a list of domain names from the the repository
function loadDomains() {
  // determine domains
  fetch("http://localhost/domain")
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Domains = yaml)
}


//------------------------------------------------------------------------------

// loadCatalog retrieves a list of components of a domain from the the repository
function loadCatalog(domain) {
  // determine domains
  fetch("http://localhost/catalog/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
    .then((list)     => model.Catalog = list)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a list of component names from the the repository
function loadComponents(domain) {
  // determine domains
  fetch("http://localhost/component/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Components = yaml)
}

//------------------------------------------------------------------------------

// loadComponents retrieves a component from the the repository
function loadComponent(domain, component) {
  // determine component
  fetch("http://localhost/component/" + domain + "/" + component)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => view.component = yaml)
}

//------------------------------------------------------------------------------

// saveComponent uploads a component to the repository
function saveComponent(domain, comp) {
  body      = jsyaml.safeDump(comp)
  component = comp.Component + " - " + comp.Version
  fetch("http://localhost/component/" + domain, {method: "POST", body: body})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// deleteComponent removes a component from the the repository
function deleteComponent(domain, component) {
  fetch("http://localhost/component/" + domain + "/" + component, {method: "DELETE"})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

//------------------------------------------------------------------------------

// loadArchitectures retrieves a list of architecture names from the the repository
function loadArchitectures(domain) {
  // determine domains
  fetch("http://localhost/architecture/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architectures = yaml)
}

//------------------------------------------------------------------------------

// loadArchitecture retrieves an  architecture from the the repository
function loadArchitecture(domain, architecture) {
  // determine domains
  fetch("http://localhost/architecture/" + domain + "/" + architecture)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architecture = yaml)
}

//------------------------------------------------------------------------------
