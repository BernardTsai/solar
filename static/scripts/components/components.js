Vue.component(
  'components',
  {
    props: ['model', 'view'],
    methods: {
      addComponent: function() {
        // find new name
        index = 1
        id    = "Comp-" + index
        do {
          id    = "Comp-" + index
          found = false
          for (idx in model.Catalog) {
            c = model.Catalog[idx]
            if (c.Component == id) {
              found = true
              break
            }
          }
          // finally found a unique name
          if (!found) {
            break
          }
          index++
        } while (true);

        // add new component
        model.Catalog.push({
          Component:     id,
          Version:       "V1.0.0",
          Configuration: "",
          Dependencies:  {}
        })

        // open editor
        view.newComponent = true
        view.component    = model.Catalog.length - 1
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
          v-if="view.component>=0"
          v-bind:model="model"
          v-bind:view="view"
          v-bind:configuration="null"
          v-bind:component="model.Catalog[view.component]">
        </componentEditor>
      </div>`
  }
)
