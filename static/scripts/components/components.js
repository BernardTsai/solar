Vue.component(
  'components',
  {
    props: ['model', 'view'],
    methods: {
      // addComponent opens editor for creating a new component
      addComponent: function() {
        // reset fields for component element editor
        this.view.ce = {
          New:            true,
          Component:      "",
          Version:        "",
          Configuration1: "",
          Dependency:     "",
          DepType:        "",
          DepComponent:   "",
          DepVersion:     "",
          Configuration2: ""
        },

        // initialise the component element of the model
        this.model.Component = {
          Component:     "unknown",
          Version:       "",
          Configuration: "",
          Dependencies:  {}
        }
      }
    },
    template: `
      <div id="components">
        <div id="container">
          <component
            v-bind:model="model"
            v-bind:view="view"
            v-bind:component="component"
            v-bind:index="index"
            v-bind:key="component.Component + ' - ' + component.Version"
            v-for="(component, index) in model.Catalog">
          </component>
          <div id="add" v-on:click="addComponent">
            <i class="fas fa-2x fa-plus-circle">
          </div>
        </div>
        <componentEditor
          v-if="model.Component"
          v-bind:model="model"
          v-bind:view="view">
        </componentEditor>
      </div>`
  }
)
