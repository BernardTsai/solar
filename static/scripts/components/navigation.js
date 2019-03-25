Vue.component(
  'navigation',
  {
    props: ['model', 'view'],
    template: `
      <div id="navigation">
        <components     v-if="view.domain!='' && view.nav == 'Components'"     :model="model" :view="view"></components>
        <architecture   v-if="view.domain!='' && view.nav == 'Architecture'"   :model="model" :view="view"></architecture>
        <solution       v-if="view.domain!='' && view.nav == 'Solution'"       :model="model" :view="view"></solution>
        <automation     v-if="view.domain!='' && view.nav == 'Automation'"     :model="model" :view="view"></automation>
        <administration v-if="                   view.nav == 'Administration'" :model="model" :view="view"></administration>
      </div>`
  }
)
