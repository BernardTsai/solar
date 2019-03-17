Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    methods: {
      // selectArchitecture pick a specific version of an architecture
       selectArchitecture: function() {
        view.architecture = document.getElementById('architectureSelector').value

        if (view.architecture != "") {
          // set solution
          view.solution = getName(view.architecture)
          view.version  = getVersion(view.architecture)

          // load architectures
          if (view.domain != "" && view.architecture != ""){
            loadArchitecture(view.domain, view.architecture)
          }
        } else {
          view.solution      = ""
          view.version       = ""
          model.Architecture = null
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
            <select id="architectureSelector" v-model="view.architecture" v-on:change="selectArchitecture">
              <option selected value="">Please select one</option>
              <option v-for="architecture in model.Architectures">{{architecture}}</option>
            </select>
          </div>
        </div>
        <div id="container" v-if="model.Architecture">
          <architectureElement
            v-bind:model="model"
            v-bind:view="view"
            v-bind:element="element"
            v-for="element in model.Architecture.Elements">
          </architectureElement>
          <div id="add" v-on:click="addElement" v-if="model.Architecture">
            <i class="fas fa-2x fa-plus-circle">
          </div>
          <architectureElementEditor
            v-bind:model="model"
            v-bind:view="view">
          </architectureElementEditor>
        </div>
      </div>`
  }
)
