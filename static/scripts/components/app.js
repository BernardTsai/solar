Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <div id="header">
          <div id="title" title="Simplified Orchestration of the Lifecycle Automation of Resources\nReference implementation - © Bernard Tsai 2018">SOLAR</div>
          <div id="selectors" v-if="view.domain!=''" >
            <div id="domain-selector">
              <strong>Domain:</strong> &nbsp; {{view.domain}} &nbsp;
            </div>
          </div>
          <div id="nav">
            <div v-if="view.domain!=''" @click="navComponents"     :class="{selected: view.nav=='Components'}"     id="Components">Catalog            <i class="fas fa-cube   text-gray-300"></i></div>
            <div v-if="view.domain!=''" @click="navArchitecture"   :class="{selected: view.nav=='Architecture'}"   id="Architecture">Architecture     <i class="fas fa-map    text-gray-300"></i></div>
            <div v-if="view.domain!=''" @click="navSolution"       :class="{selected: view.nav=='Solution'}"       id="Solution">Solution             <i class="fas fa-cubes  text-gray-300"></i></div>
            <div v-if="view.domain!=''" @click="navAutomation"     :class="{selected: view.nav=='Automation'}"     id="Automation">Automation         <i class="fas fa-cogs   text-gray-300"></i></div>
            <div                        @click="navAdministration" :class="{selected: view.nav=='Administration'}" id="Administration">Administration <i class="fas fa-wrench text-gray-300"></i></div>
          </div>
        </div>
        <navigation v-bind:model="model" v-bind:view="view"></navigation>
      </div>`
  }
)

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

//------------------------------------------------------------------------------

function navSolution()     {
  view.nav = "Solution"


  model.SolElement = null // no element to be displayed
  model.Graph      = null // reset graph

  // check if the solution selection needs to be updated
  if (!model.Solutions) {
    loadSolutions(view.domain)
    .then(() => {
      // load solution if necessary
      if (model.Solutions.includes(view.solution) && !model.Solution) {
        loadSolution(view.domain, view.solution)
      }
    })
  } else {
    // load solution if necessary
    if (model.Solutions.includes(view.solution) && !model.Solution) {
      loadSolution(view.domain, view.solution)
    }
  }
}

//------------------------------------------------------------------------------

function navAutomation()   {
  view.nav = "Automation";

  view.automation.solution = ""
  view.automation.element  = ""
  view.automation.cluster  = ""
  view.automation.instance = ""

  model.Tasks = []
  model.Trace = null
}

//------------------------------------------------------------------------------

function navAdministration() {
  view.nav = "Administration";
}

//------------------------------------------------------------------------------
