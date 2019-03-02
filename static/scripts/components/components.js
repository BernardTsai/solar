Vue.component(
  'components',
  {
    props: ['model', 'view'],
    template: `
      <div id="components">
        <div id="container">
          <component
            v-bind:model="model"
            v-bind:view="view"
            v-bind:component="component"
            v-for="component in model.Components">
          </component>
          <component
            v-bind:model="model"
            v-bind:view="view"
            v-bind:component="''">
          </component>
        </div>
        <componentEditor
          v-id="view.component!=''"
          v-bind:model="model"
          v-bind:view="view"
          v-bind:component="view.component">
        </componentEditor>
      </div>`
  }
)
