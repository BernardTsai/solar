var model = {
  Domains:       [],
  Catalog:       [],
  Components:    [],
  Architectures: {},
  Architecture:  null,
  Solution:      {},
  Tasks:         [],
  Task:          {}
};

function loadDomains() {
  // determine domains
  fetch("http://localhost/domain")
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Domains = yaml)
}

function loadCatalog(domain) {
  // determine domains
  fetch("http://localhost/catalog/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => yaml.sort((a,b) => (a.Component > b.Component) ? 1 : (a.Component < b.Component) ? -1 : 0))
    .then((list)     => model.Catalog = list)
}

function loadComponents(domain) {
  // determine domains
  fetch("http://localhost/component/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Components = yaml)
}

function loadComponent(domain, component) {
  // determine component
  fetch("http://localhost/component/" + domain + "/" + component)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => view.component = yaml)
}

function saveComponent(domain, comp) {
  body      = jsyaml.safeDump(comp)
  component = comp.Component + " - " + comp.Version
  fetch("http://localhost/component/" + domain, {method: "POST", body: body})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

function deleteComponent(domain, component) {
  fetch("http://localhost/component/" + domain + "/" + component, {method: "DELETE"})
    .then((response) => response.text())
    .then((text)     => loadCatalog(domain))
}

function loadArchitectures(domain) {
  // determine domains
  fetch("http://localhost/architecture/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architectures = yaml)
}

function loadArchitecture(domain, architecture) {
  // determine domains
  fetch("http://localhost/architecture/" + domain + "/" + architecture)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Architecture = yaml)
}
