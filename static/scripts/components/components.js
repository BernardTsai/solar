Vue.component(
  'components',
  {
    props: ['model', 'view'],
    methods: {
      // editComponent opens the editor for the selected component
      editComponent: function(component) {
        // reset fields for component editor
        this.view.ce = {
          New:        false,
          Component:  component.Component,
          Version:    component.Version,
          Dependency: "",
          ConfTitle:  ""
        },

        // initialise the component of the model
        this.model.Component = component
      },
      // addComponent opens editor for creating a new component
      addComponent: function() {
        // reset fields for component element editor
        this.view.ce = {
          New:        true,
          Component:  "unknown",
          Version:    "V1.0.0",
          Dependency: "",
          ConfTitle:  ""
        },

        // initialise the component element of the model
        this.model.Component = {
          Component:     "unknown",
          Version:       "V1.0.0",
          Configuration: "",
          Dependencies:  {}
        }
      }
    },
    template: `
      <div id="components">
        <div id="comps">
          <div class="header">
            <h3>Catalog:</h3>
          </div>

          <table class="components">
            <thead>
              <tr>
                <th>Component</th>
                <th>Version</th>
                <th class="center" @click="addComponent"><i class="fas fa-plus-circle"></i></th>
            </thead>
            <tbody>
              <tr v-for="comp in model.Catalog">
                <td>{{comp.Component}}</td>
                <td>{{comp.Version}}</td>
                <td @click="editComponent(comp)">
                  <i class="fas fa-edit"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <compEditor
          v-if="model.Component"
          v-bind:model="model"
          v-bind:view="view">
        </compEditor>
      </div>`
  }
)
