Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    methods: {
      addElement: function() {
        view.newElement = true
        view.element    = {
          Element:       "unknown",
          Component:     "",
          Configuration: "",
          Clusters:      []
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
          <architectureElementWizard
            v-bind:model="model"
            v-bind:view="view"
            v-bind:element="{}">
          </architectureElementWizard>
        </div>
      </div>`
  }
)
