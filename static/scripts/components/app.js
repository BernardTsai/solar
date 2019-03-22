Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <div id="header">
          <div id="title" title="Simplified Orchestration of the Lifecycle Automation of Resources\nReference implementation - © Bernard Tsai 2018">SOLAR</div>
          <div id="selectors">
            <div id="domain-selector">
              <strong>Domain:</strong>
              <select id="domainSelector" v-model="view.domain" @change="selectDomain">
                <option selected value="">Please select one</option>
                <option v-for="domain in model.Domains">{{domain}}</option>
              </select>
            </div>
          </div>
          <div id="nav" v-if="view.domain!=''">
            <div v-on:click="navComponents"      :class="{selected: view.nav=='Components'}"     id="Components">Catalog            <i class="fas fa-cube   text-gray-300"></i></div>
            <div v-on:click="navArchitecture"    :class="{selected: view.nav=='Architecture'}"   id="Architecture">Architecture     <i class="fas fa-map    text-gray-300"></i></div>
            <div v-on:click="navSolution"        :class="{selected: view.nav=='Solution'}"       id="Solution">Solution             <i class="fas fa-cubes  text-gray-300"></i></div>
            <div v-on:click="navAutomation"      :class="{selected: view.nav=='Automation'}"     id="Automation">Automation         <i class="fas fa-cogs   text-gray-300"></i></div>
          </div>
          <div id="nav" v-if="view.domain==''">
            <div v-on:click="navAdministration"  :class="{selected: view.nav=='Administration'}" id="Administration">Administration <i class="fas fa-wrench text-gray-300"></i></div>
          </div>
        </div>
        <navigation v-bind:model="model" v-bind:view="view"></navigation>
      </div>`
  }
)

//------------------------------------------------------------------------------

// selectArchitecture pick a specific version of an architecture
function selectDomain() {
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

  // load solutions
  if (view.domain == ""){
    model.Solutions = []
    model.Solution  = null
  } else {
    loadSolutions(view.domain)
    model.Solution  = null
  }
  view.solution = ""

  // activate administration if no domain has been selected
  if (view.domain == "") {
    navAdministration()
  } else if (view.nav == "" || view.nav == "Administration") {
    view.nav = "Components"
  }
}

//------------------------------------------------------------------------------

function navComponents()   {
  view.nav = "Components"

  model.Component = null
}

//------------------------------------------------------------------------------

function navArchitecture() {
  view.nav = "Architecture";

  model.ArchElement = null;
  model.Graph       = null;
}
function navSolution()     {
  view.nav = "Solution";

  model.SolElement = null;
  model.Graph      = null;
}
function navAutomation()   {
  view.nav = "Automation";

  view.automation.solution = ""
  view.automation.element  = ""
  view.automation.cluster  = ""
  view.automation.instance = ""

  model.Tasks = []
  model.Trace = null
}
function navAdministration() {
  view.nav = "Administration";
}

//------------------------------------------------------------------------------
