var model = {
  Domains:       [],
  Components:    [],
  Architectures: {},
  Architecture:  {},
  Solution:      {},
  Tasks:         [],
  Task:          {}
};

function loadModel() {
  // determine domains
  fetch("http://localhost/domain")
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Domains = yaml)
}

function loadComponents(domain) {
  // determine domains
  fetch("http://localhost/component/" + domain)
    .then((response) => response.text())
    .then((text)     => jsyaml.safeLoad(text))
    .then((yaml)     => model.Components = yaml)
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
