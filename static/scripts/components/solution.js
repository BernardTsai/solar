Vue.component(
  'solution',
  {
    props: ['model', 'view'],
    methods: {
      // viewElement displays an element in the editor
      viewElement: function(element) {
        // initialise the solution element of the model
        view.Element  = element.Element
        model.Element = element

        this.$forceUpdate()
      },
      // hidelement hides the editor
      hideElement: function(element) {
        // reset the solution element of the model
        view.Element  = ""
        model.Element = null

        this.$forceUpdate()
      },
      // selectSolution pick a specific solution
      selectSolution: function() {
        // view.solution = document.getElementById('solutionSelector').value

        // load solution
        if (this.view.domain != "" && this.view.solution != ""){
          this.view.graph.viewElement = this.viewElement
          loadSolution(this.view.domain, this.view.solution)
        } else {
          this.model.Solution = null
          this.model.Graph    = null
        }
        // reset element
        this.model.Element = null

        this.$forceUpdate()
      }
    },
    computed: {
      // graph creates the solution graph
      graph: function() {
        if (this.view.solution != '') {
          sg = new SolutionGraph(this.model, this.view, this.view.domain, this.view.solution)
        }
        return ""
      }
    },
    template: `
      <div id="solution" v-if="view.nav=='Solution'">
        {{graph}}

        <div id="selector">
          <div id="solution-selector">
            <strong>Solution:</strong>
            <select id="solutionSelector" v-model="view.solution" v-on:change="selectSolution">
              <option selected value="">Please select one</option>
              <option v-for="solution in model.Solutions">{{solution}}</option>
            </select>
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

        <div id="container" v-if="model.Graph">
          <svg id="canvas" v-bind:style="{ width: model.Graph.Width + 'px', height: model.Graph.Height + 'px'}">
            <edge
              :model="model"
              :view="view"
              :edge="edge"
              v-for="edge in model.Graph.Edges">
            </edge>
            <node
              :model="model"
              :view="view"
              :node="node"
              v-for="node in model.Graph.Nodes">
            </node>
          </svg>
        </div>

        <solEditor v-if="model.Element" :model="model" :view="view" :element="model.Element"/>
      </div>`
  }
)

//------------------------------------------------------------------------------
