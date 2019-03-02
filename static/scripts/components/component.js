Vue.component(
  'component',
  {
    props: ['model', 'view', 'component'],
    computed: {
      label: function() {
        if (this.component == "") {
          return "new component"
        }
        return this.component
      },
      name: function() {
        if (this.component == "") {
          return "+"
        }
        component = this.component
        version   = component.split(" - ").slice(-1)[0]
        name      = component.substring(0, component.length - 3 - version.length)

        return name
      },
      version: function() {
        if (this.component == "") {
          return "Vx.y.z"
        }
        component = this.component
        version   = component.split(" - ").slice(-1)[0]
        name      = component.substring(0, component.length - 3 - version.length)

        return version
      }
    },
    template: `
      <div  class="component" v-bind:title="label" v-bind:class="{new: name=='+'}" v-on:click="editComponent(label, $event)">
        <div class="label">
          <div class="name">{{name}}</div>
          <div class="version">{{version}}</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)

function editComponent(message, event) {
  console.log(message)
  console.log(event)
}
