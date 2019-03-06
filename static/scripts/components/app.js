Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <div id="header">
          <div id="title" title="Simplified Orchestration of the Lifecycle Automation of Resources">SOLAR</div>
          <div id="selectors">
            <div id="domain-selector">
              <strong>Domain:</strong>
              <select id="domainSelector" v-on:change="selectDomain">
                <option disabled selected value="">Please select one</option>
                <option v-for="domain in model.Domains">{{domain}}</option>
              </select>
            </div>
            <div id="architecture-selector" v-if="view.nav=='Architecture'">
              <strong>Architecture:</strong>
              <select id="architectureSelector" v-model="view.architecture" v-on:change="selectArchitecture">
                <option disabled selected value="">Please select one</option>
                <option v-for="architecture in model.Architectures">{{architecture}}</option>
              </select>
            </div>
          </div>
          <div id="nav" v-if="view.domain!=''">
            <div v-on:click="navComponents"   :class="{selected: view.nav=='Components'}"   id="Components">Components     <i class="fas fa-cube  text-gray-300"></i></div>
            <div v-on:click="navArchitecture" :class="{selected: view.nav=='Architecture'}" id="Architecture">Architecture <i class="fas fa-map   text-gray-300"></i></div>
            <div v-on:click="navSolution"     :class="{selected: view.nav=='Solution'}"     id="Solution">Solution         <i class="fas fa-cubes text-gray-300"></i></div>
            <div v-on:click="navAutomation"   :class="{selected: view.nav=='Automation'}"   id="Automation">Automation     <i class="fas fa-cogs  text-gray-300"></i></div>
          </div>
        </div>
        <navigation v-bind:model="model" v-bind:view="view"></navigation>
      </div>`
  }
)

function selectDomain() {
  view.domain = document.getElementById('domainSelector').value

  // load catalog
  if (view.domain == ""){
    model.Catalog = []
  } else {
    loadCatalog(view.domain)
  }

  // load components
  if (view.domain == ""){
    model.Components = []
  } else {
    loadComponents(view.domain)
  }

  // load architectures
  if (view.domain == ""){
    model.Architectures = []
    model.Architecture  = null
  } else {
    loadArchitectures(view.domain)
    model.Architecture  = null
  }
  view.architecture = ""
}

function selectArchitecture() {
  view.architecture = document.getElementById('architectureSelector').value

  // set solution
  view.solution = getName(view.architecture)
  view.version  = getVersion(view.architecture)

  // load architectures
  if (view.domain != "" && view.architecture != ""){
    loadArchitecture(view.domain, view.architecture)
  }
}

function navComponents()   { view.nav = "Components";   }
function navArchitecture() { view.nav = "Architecture"; }
function navSolution()     { view.nav = "Solution"; }
function navAutomation()   { view.nav = "Automation"; }
