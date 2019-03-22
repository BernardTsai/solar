Vue.component(
  'navigation',
  {
    props: ['model', 'view'],
    template: `
      <div id="navigation">
        <components     v-if="view.domain!='' && view.nav == 'Components'"     v-bind:model="model" v-bind:view="view"></components>
        <architecture   v-if="view.domain!='' && view.nav == 'Architecture'"   v-bind:model="model" v-bind:view="view"></architecture>
        <solution       v-if="view.domain!='' && view.nav == 'Solution'"       v-bind:model="model" v-bind:view="view"></solution>
        <automation     v-if="view.domain!='' && view.nav == 'Automation'"     v-bind:model="model" v-bind:view="view"></automation>
        <administration v-if="view.domain=='' && view.nav == 'Administration'" v-bind:model="model" v-bind:view="view"></administration>
      </div>`
  }
)
