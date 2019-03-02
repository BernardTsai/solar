Vue.component(
  'architectureElement',
  {
    props: ['model', 'view', 'element'],
    computed: {
      label: function() {
        return this.element.Element
      }
    },
    template: `
      <div  class="architectureElement" v-bind:title="label">
        <div class="label">
          <div class="name">{{label}}</div>
          <div class="version">&nbsp;</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)
