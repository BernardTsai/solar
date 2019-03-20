Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    methods: {
      // selectArchitecture pick a specific version of an architecture
       selectArchitecture: function() {
        this.view.architecture = document.getElementById('architectureSelector').value

        if (view.architecture != "") {
          // set solution
          this.view.solution = getName(this.view.architecture)
          this.view.version  = getVersion(this.view.architecture)

          // load architectures
          if (this.view.domain != "" && this.view.architecture != ""){
            loadArchitecture(this.view.domain, this.view.architecture)
            ag = new ArchitectureGraph(this.model, this.view, this.view.domain, this.view.solution, this.view.version)
          }
        } else {
          this.view.solution      = ""
          this.view.version       = ""
          this.model.Architecture = null
        }
      },
      // addElement adds a new architecture element
      addElement: function() {
        // reset fields for architecture element editor
        this.view.ae = {
          New:            true,
          Element:        "",
          Component:      "",
          Configuration1: "",
          Cluster:        "",
          State:          "",
          Min:            "",
          Max:            "",
          Size:           "",
          Configuration2: "",
          Relationship:   "",
          Dependency:     "",
          DepType:        "",
          RelElement:     "",
          Configuration3: ""
        },

        // initialise the architecture element of the model
        this.model.ArchElement = {
          Element:       "unknown",
          Component:     "",
          Configuration: "",
          Clusters:      {}
        }
      }
    },
    template: `
      <div id="architecture" v-if="view.nav=='Architecture'">
        <div id="selector">
          <div id="architecture-selector">
            <strong>Architecture:</strong>
            <select id="architectureSelector" v-model="view.architecture" @change="selectArchitecture">
              <option selected value="">Please select one</option>
              <option v-for="architecture in model.Architectures">{{architecture}}</option>
            </select>
          </div>

          <button class="modal-default-button">
            Duplicate <i class="fas fa-copy">
          </button>
          <button class="modal-default-button">
            Delete <i class="fas fa-times-circle">
          </button>
          <button class="modal-default-button">
            Update <i class="fas fa-plus-circle">
          </button>
          <button class="modal-default-button" >
            Create <i class="fas fa-plus-circle">
          </button>

        </div>

        <div id="comps">
          <div class="header">
            <h3>Catalog:</h3>
          </div>

          <table class="components">
            <thead>
              <tr>
                <th>Component</th>
                <th>Version</th>
                <th></th>
            </thead>
            <tbody>
              <tr v-for="comp in model.Catalog">
                <td>{{comp.Component}}</td>
                <td>{{comp.Version}}</td>
                <td>
                  <i class="fas fa-cube"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div id="container" v-if="model.Architecture && model.Graph">

          <svg id="canvas" v-bind:style="{ width: model.Graph.Width + 'px', height: model.Graph.Height + 'px'}">
            <edge
              v-bind:model="model"
              v-bind:view="view"
              v-bind:edge="edge"
              v-for="edge in model.Graph.Edges">
            </edge>
            <node
              v-bind:model="model"
              v-bind:view="view"
              v-bind:node="node"
              v-for="node in model.Graph.Nodes">
            </node>
          </svg>

          <!-- architectureElement
            v-bind:model="model"
            v-bind:view="view"
            v-bind:element="element"
            v-for="element in model.Architecture.Elements">
          </architectureElement>

          <div id="add" v-on:click="addElement" v-if="model.Architecture">
            <i class="fas fa-2x fa-plus-circle">
          </div -->

        </div>

        <architectureElementEditor
          v-bind:model="model"
          v-bind:view="view">
        </architectureElementEditor>
      </div>`
  }
)
