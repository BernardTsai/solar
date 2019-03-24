Vue.component(
  'solution',
  {
    props: ['model', 'view'],
    methods: {
      // showTasks switches to the corresponding automation tasks
      showTasks: function() {
        model.Tasks = []
        model.Trace = null

        view.automation.solution = view.solution
        view.automation.element  = ""
        view.automation.cluster  = ""
        view.automation.instance = ""

        loadTasks(
          view.domain,
          view.automation.solution,
          view.automation.element,
          view.automation.cluster,
          view.automation.instance)

        view.nav = "Automation"
      },
      // refresh reloads the solution
      refresh: function() {
        if (this.view.domain != "" && this.view.solution != ""){
          loadSolution(this.view.domain, this.view.solution)
          .then(() => {
            // update selected element if necessary
            if (this.model.SolElement)
            {
              element = this.model.SolElement.Element

              this.model.SolElement = this.model.Solution.Elements[element]
            }
          })
        }
        // reset element
        // this.model.SolElement = null
      },
      // viewElement displays an element in the editor
      viewElement: function(element) {
        // initialise the solution element of the model
        this.model.SolElement = element
      },
      // viewNode displays a nodes element in the editor
      viewNode: function(node) {
        // initialise the solution element of the model
        this.model.SolElement = node.Element
      },
      // hidelement hides the editor
      hideElement: function() {
        // reset the solution element of the model
        this.model.SolElement = null
      },
      // selectSolution pick a specific solution
      selectSolution: function() {
        // load solution
        if (this.view.domain != "" && this.view.solution != ""){
          loadSolution(this.view.domain, this.view.solution)
        } else {
          this.model.Solution = null
          this.model.Graph    = null
        }
        // reset element
        this.model.SolElement = null
      },
      // graph creates the solution graph
      graph: function() {
        return new SolutionGraph(this.model, this.view, this.view.domain, this.view.solution)
      }
    },
    template: `
      <div id="solution" v-if="view.nav=='Solution'">

        <div id="selector">
          <div id="solution-selector">
            <strong>Solution:</strong>
            <select id="solutionSelector" v-model="view.solution" @change="selectSolution">
              <option selected value="">Please select one</option>
              <option v-for="solution in model.Solutions">{{solution}}</option>
            </select>
          </div>

          <div class="buttons">
            <button class="action" v-if="view.solution!=''" @click="refresh()">
              Refresh <i class="fas fa-recycle">
            </button>
            <button class="action" v-if="view.solution!=''" @click="showTasks()">
              Tasks <i class="fas fa-cogs">
            </button>
          </div>
        </div>

        <div id="elements" v-if="model.Solution">
          <div class="header">
            <h3>Elements:</h3>
          </div>

          <table class="elements">
            <thead>
              <tr>
                <th>Element</th>
                <th>Type</th>
                <th @click="hideElement()">
                  <i class="fas fa-eye-slash"></i>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="element in model.Solution.Elements">
                <td>{{element.Element}}</td>
                <td>{{element.Component}}</td>
                <td @click="viewElement(element)">
                  <i class="fas fa-eye"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div id="container" v-if="model.Solution">
          <graph :model="model" :view="view" :graph="graph()" @node-selected="viewNode"/>
        </div>

        <solEditor v-if="model.SolElement" :model="model" :view="view" :element="model.SolElement"/>
      </div>`
  }
)

//------------------------------------------------------------------------------
