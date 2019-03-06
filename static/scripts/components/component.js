Vue.component(
  'component',
  {
    props: ['model', 'view', 'component', 'index'],
    methods: {
      editComponent: function() {
        // reset fields for component editor
        this.view.ce = {
          New:            false,
          Component:      this.component.Component,
          Version:        this.component.Version,
          Configuration1: this.component.Configuration,
          Dependency:     "",
          DepType:        "",
          DepComponent:   "",
          DepVersion:     "",
          Configuration2: ""
        },

        // initialise the component of the model
        this.model.Component = this.component
      }
    },
    computed: {
      label: function() {
        return this.name + " - " + this.version
      },
      name: function() {
        if (this.component == null) {
          return "Unknown"
        }
        return this.component.Component
      },
      version: function() {
        if (this.component == null) {
          return "Vx.y.z"
        }
        return this.component.Version
      }
    },
    template: `
      <div  class="component"
        v-bind:title="label"
        v-bind:class="{new: name=='Unknown'}"
        v-on:click="editComponent()">
        <div class="label">
          <div class="name">{{name}}</div>
          <div class="version">{{version}}</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)
