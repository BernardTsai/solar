Vue.component(
  'componentEditor',
  {
    props: ['model', 'view', 'component','configuration'],
    computed: {
      components: function() {
        var map = {};
        for (var index in model.Catalog) {
          component = model.Catalog[index]
          version   = component.Version
          name      = component.Component

          c = map[name]
          if (c == null) {
            map[name] = name
          }
        }
        var result = []
        for (var component in map) {
          result.push(component)
        }
        return result
      }
    },
    methods: {
      newVersion: function(oldVersion) {
        parts = oldVersion.substring(1).split(".")

        nextPatch = parseInt(parts[2]) + 1

        return "V" + parts[0] + "." + parts[1] + "." + nextPatch

      },
      editConfiguration: function(element) {
        this.configuration = element
      },
      closeConfigurationEditor: function() {
        this.configuration = null
      },
      cancelDialog: function() {
        if (this.view.newComponent) {
            this.model.Catalog.splice(-1, 1)
        }
        this.view.component = -1
      },
      saveDialog: function() {
        view.component = -1
        saveComponent(this.view.domain, this.component)
      },
      deleteDialog: function() {
        view.component = -1
        deleteComponent(this.view.domain, this.component.Component + " - " + this.component.Version)
      },
      duplicateDialog: function() {
        // deep copy component
        comp = JSON.parse(JSON.stringify(this.component))

        // update version

        comp.Version = this.newVersion(this.component.Version)

        // add to catalog
        this.model.Catalog.push( comp )

        // open editor
        view.newComponent = true
        view.component    = model.Catalog.length - 1
      },
      addDependency: function() {
        // find new name
        index = 1
        id    = "Dep-" + index
        do {
          id    = "Dep-" + index
          found = false
          for (idx in this.component.Dependencies) {
            d = this.component.Dependencies[idx]
            if (d.Dependency == id) {
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

        // add dependency
        this.component.Dependencies[id] = {
          Dependency:     id,
          Type:           "service",
          Component:      "",
          Version:        "",
          Configuration:  ""
        }
        this.$forceUpdate();
      },
      deleteDependency: function (dependency) {
        delete this.component.Dependencies[dependency]
        this.$forceUpdate();
      },
      dependencyChanged: function() {
        this.$forceUpdate();
      }
    },
    template: `
      <div class="modal-mask">
        <div class="modal-wrapper">

          <div class="modal-editor" v-if="configuration!=null">
            <div class="modal-header">
              <h3>Component: {{component.Component}} - {{component.Version}}</h3>
              <div class="close" v-on:click="cancelDialog()"><i class="far fa-times-circle"></i></div>
            </div>
            <div class="modal-body">
              <textarea autofocus cols=80 rows=10 v-model="configuration.Configuration"></textarea>
            </div>
            <div class="modal-footer">
              &nbsp;
              <button class="modal-default-button" v-on:click="closeConfigurationEditor()">
                <strong>OK</strong> <i class="fas fa-check">
              </button>
            </div>
          </div>

          <div class="modal-container" v-if="configuration==null">

            <div class="modal-header">
              <h3>Component: {{component.Component}} - {{component.Version}}</h3>
              <div class="close" v-on:click="cancelDialog()"><i class="far fa-times-circle"></i></div>
            </div>

            <div class="modal-body">
              <table style="width: 100%">
                <col width="1*">
                <col width="1*">
                <col width="1*">
                <col width="1*">
                <col width="80%">
                <col width="0*">
                <tr>
                  <td><strong>Component:</strong></td>
                  <td colspan=5>
                    <input name="component" class="long" type="text" v-model="component.Component" :disabled="!view.newComponent"></input>
                  </td>
                </tr>
                <tr>
                  <td><strong>Version:</strong></td>
                  <td colspan=5>
                    <input name="version" class="long" type="text" v-model="component.Version" :disabled="!view.newComponent"></input>
                  </td>
                </tr>
                <tr>
                  <td><strong>Configuration:</strong></td>
                  <td colspan=5>
                    <input class="long" type="text" v-model="component.Configuration" v-on:click="editConfiguration(component)" :disabled="!view.newComponent"></input>
                  </td>
                </tr>
                <tr>
                  <td colspan=6>&nbsp;</td>
                </tr>
                <tr>
                  <td><strong>Dependency:</strong></td>
                  <td><strong>Dep.&nbsp;Type:</strong></td>
                  <td><strong>Component:</strong></td>
                  <td><strong>Version:</strong></td>
                  <td><strong>Configuration:</strong></td>
                  <td v-on:click="addDependency()"><strong><i class="fas fa-plus-circle" v-if="view.newComponent"></i></strong></td>
                </tr>
                <tr v-for="dependency in component.Dependencies" :key="dependency.Dependency">
                  <td>
                    <input type="text" v-model="dependency.Dependency" :disabled="!view.newComponent"></input>
                  </td>
                  <td>
                    <select v-model="dependency.Type" :disabled="!view.newComponent">
                      <option disabled value="">please select</option>
                      <option>context</option>
                      <option>service</option>
                    </select>
                  </td>
                  <td>
                    <select v-model="dependency.Component" v-on:change="dependencyChanged()" :disabled="!view.newComponent">
                      <option disabled value="">component</option>
                      <option v-for="c in components">{{c}}</option>
                    </select>
                  </td>
                  <td>
                    <select v-model="dependency.Version" :disabled="!view.newComponent">
                      <option disabled value="">Vx.y.z</option>
                      <option v-for="c in model.Catalog" v-if="c.Component == dependency.Component">{{c.Version}}</option>
                    </select>
                  </td>
                  <td>
                    <input type="text" v-model="dependency.Configuration" v-on:click="editConfiguration(dependency)" :disabled="!view.newComponent"></input>
                  </td>
                  <td v-on:click="deleteDependency(dependency.Dependency)"><i class="fas fa-minus-circle" v-if="view.newComponent"></i></td>
                </tr>
              </table>
            </div>

            <div class="modal-footer">
              &nbsp;
              <button class="modal-default-button" v-on:click="duplicateDialog()" v-if="!view.newComponent">
                <strong>Duplicate</strong> <i class="fas fa-copy">
              </button>
              <button class="modal-default-button" v-on:click="deleteDialog()" v-if="!view.newComponent">
                <strong>Delete</strong> <i class="fas fa-trash-alt">
              </button>
              <button class="modal-default-button" v-on:click="saveDialog()" v-if="view.newComponent">
                <strong>Save</strong> <i class="fas fa-save">
              </button>
            </div>
          </div>
        </div>
      </div>`
  }
)
