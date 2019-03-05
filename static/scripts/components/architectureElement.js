Vue.component(
  'architectureElement',
  {
    props: ['model', 'view', 'element'],
    methods: {
      editElement: function() {
        view.Element    = this.element
        view.newElement = false
      }
    },
    computed: {
      label: function() {
        return this.element.Element
      }
    },
    template: `
      <div  class="architectureElement" v-bind:title="label" v-on:click="editElement()">
        <div class="label">
          <div class="name">{{label}}</div>
          <div class="version">&nbsp;</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)
