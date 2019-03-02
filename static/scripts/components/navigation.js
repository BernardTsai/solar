Vue.component(
  'navigation',
  {
    props: ['model', 'view'],
    template: `
      <div id="navigation" v-if="view.domain!=''">
        <components   v-if="view.nav == 'Components'"   v-bind:model="model" v-bind:view="view"></components>
        <architecture v-if="view.nav == 'Architecture'" v-bind:model="model" v-bind:view="view"></architecture>
        <solution     v-if="view.nav == 'Solution'"     v-bind:model="model" v-bind:view="view"></solution>
        <automation   v-if="view.nav == 'Automation'"   v-bind:model="model" v-bind:view="view"></automation>
      </div>`
  }
)
