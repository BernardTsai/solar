Vue.component(
  'component',
  {
    props: ['model', 'view', 'component', 'index'],
    methods: {
      editComponent: function(index) {
        view.newComponent = false
        view.component    = index
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
        v-on:click="editComponent(index)">
        <div class="label">
          <div class="name">{{name}}</div>
          <div class="version">{{version}}</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)
