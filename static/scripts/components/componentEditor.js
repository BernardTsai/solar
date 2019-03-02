Vue.component(
  'componentEditor',
  {
    props: ['model', 'view', 'component'],
    computed: {
      label: function() {
        if (this.component == "") {
          return "new component"
        }
        return this.component
      },
      name: function() {
        if (this.component == "") {
          return "+"
        }
        component = this.component
        version   = component.split(" - ").slice(-1)[0]
        name      = component.substring(0, component.length - 3 - version.length)

        return name
      },
      version: function() {
        if (this.component == "") {
          return "Vx.y.z"
        }
        component = this.component
        version   = component.split(" - ").slice(-1)[0]
        name      = component.substring(0, component.length - 3 - version.length)

        return version
      }
    },
    template: `
      <div class="modal-mask" v-if="view.component!=''">
        <div class="modal-wrapper">
          <div class="modal-container">

            <div class="modal-header">
              <slot name="header">
                default header
              </slot>
            </div>

            <div class="modal-body">
              <slot name="body">
                default body
              </slot>
            </div>

            <div class="modal-footer">
              <slot name="footer">
                default footer
                <button class="modal-default-button" @click="$emit('close')">
                  OK
                </button>
              </slot>
            </div>
          </div>
        </div>
      </div>`
  }
)
