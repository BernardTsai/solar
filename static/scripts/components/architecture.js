Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    methods: {
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
