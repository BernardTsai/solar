Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    template: `
      <div id="architecture">
        <div id="container" v-if="model.Architecture">
          <architectureElement
            v-bind:model="model"
            v-bind:view="view"
            v-bind:element="element"
            v-for="element in model.Architecture.Elements">
          </architectureElement>
        </div>
      </div>`
  }
)
