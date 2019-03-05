Vue.component(
  'architectureElement',
  {
    props: ['model', 'view', 'element'],
    methods: {
      editElement: function() {
        // reset fields for architecture element editor
        this.view.ae = {
          New:            false,
          Element:        this.element.Element,
          Component:      this.element.Component,
          Configuration1: this.element.Configuration,
          Cluster:        "",
          State:          "",
          Min:            "",
          Max:            "",
          Size:           "",
          Configuration2: "",
          Relationship:   "",
          Dependency:     "",
          DepType:        "",
          RelElement:     "",
          Configuration3: ""
        },

        // initialise the architecture element of the model
        this.model.ArchElement = this.element
      }
    },
    template: `
      <div  class="architectureElement" v-bind:title="element.Element" v-on:click="editElement()">
        <div class="label">
          <div class="name">{{element.Element}}</div>
          <div class="version">&nbsp;</div>
        </div>
        <div class="logo"><i class="fas fa-2x fa-cube"></i></div>
      </div>`
  }
)
