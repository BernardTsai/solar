Vue.component(
  'componentEditor',
  {
    props: ['model', 'view'],
    methods: {
      // addDependency adds a new dependency to a component
      addDependency: function() {
        // ask for name of new dependency and add it with initial values
        name = prompt("Name of the new dependency:")
        if (name && name != "") {
          // add dependency
          this.model.Component.Dependencies[name] = {
            Dependency:    name,
            Type:          "",
            Component:     "",
            Version:       "",
            Configuration: ""
          }

          // update dialog
          this.view.ce.Dependency     = name
          this.view.ce.DepType        = ""
          this.view.ce.DepComponent   = ""
          this.view.ce.DepVersion     = ""
          this.view.ce.Configuration2 = ""
        }
      },
      // delDependency removes an existing dependency from a component
      delDependency: function() {
        // check if dependency has been selected
        if (this.view.ce.Dependency=="") {
          return
        }

        // remove dependency
        delete this.model.Component.Dependencies[this.view.ce.Dependency]

        // update the fields
        this.view.ce.Dependency     = ""
        this.view.ce.DepType        = ""
        this.view.ce.DepComponent   = ""
        this.view.ce.DepVersion     = ""
        this.view.ce.Configuration2 = ""
      },
      // cancelDialog closes the dialog
      cancelDialog: function() {
        this.model.Component = null
        this.view.ce.New     = false
      },
      // createComponent adds component to the catalog
      createComponent: function() {
        saveComponent(this.view.domain, this.model.Component)
        this.model.Component = null
      },
      // updateComponent updates an existing component of the catalog
      updateComponent: function() {
        deleteComponent(this.view.domain, this.model.Component.Component + " - " + this.model.Component.Version)
        saveComponent(this.view.domain, this.model.Component)
        this.model.Component = null
      },
      // deleteComponent removes an existing component from the catalog
      deleteComponent: function() {
        deleteComponent(this.view.domain, this.model.Component.Component + " - " + this.model.Component.Version)
        this.model.Component = null
      },
      // duplicateComponent makes a copy of an existing component of the catalog
      duplicateComponent: function() {
        // update version
        version = this.newVersion(this.model.Component.Version)

        // mark component as new
        this.view.ce.New = true

        // update in edtior
        this.view.ce.Version = version

        // update model
        this.model.Component.Version = version
      },
      // newVersion determines a new version number by incrementing the patch
      newVersion: function(oldVersion) {
        parts = oldVersion.substring(1).split(".")

        nextPatch = parseInt(parts[2]) + 1

        return "V" + parts[0] + "." + parts[1] + "." + nextPatch
      },
      // components lists the name of available components in the Catalog
      components: function() {
        var map    = {};
        var result = []

        // loop over all available components
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
          name      = component.Component

          // add new components to result
          if (!map[name]) {
            map[name] = name
            result.push(name)
          }
        }
        return result
      },
      // versions lists the name of corresponding versions in the Catalog
      versions: function() {
        var result = []

        // loop over all available components
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]

          // add new version to result
          if (component.Component == view.ce.DepComponent) {
            result.push(component.Version)
          }
        }
        return result
      },
      // dependencies lists the name of dependencies of a component
      dependencies: function() {
        var result = []

        // loop over all dependencies of the component
        for (var index in this.model.Component.Dependencies) {
          dependency = this.model.Component.Dependencies[index]

          result.push(dependency.Dependency)
        }
        return result
      },
      // componentChanged  handles changes of the component information
      componentChanged: function() {
        this.model.Component.Component = this.view.ce.Component
      },
      // versionChanged handles changes of the version information
      versionChanged: function() {
        this.model.Component.Version = this.view.ce.Version
      },
      // configuration1Changed  handles changes of the component configuration information
      configuration1Changed: function() {
        this.model.Component.Configuration = this.view.ce.Configuration1
      },
      // configuration1Focus expands the configuration editor of the component configuration
      configuration1Focus: function() {
        document.getElementById("ceConfiguration1").setAttribute("rows", 10)
      },
      // configuration1Blur schrinks the configuration editor of the component configuration
      configuration1Blur: function() {
        document.getElementById("ceConfiguration1").setAttribute("rows", 2)
      },

      // dependencyChanged handles changes of the dependency information
      dependencyChanged: function() {
        // determine dependency
        dependency = this.model.Component.Dependencies[this.view.ce.Dependency]

        // update the fields
        this.view.ce.Dependency     = dependency.Dependency
        this.view.ce.DepType        = dependency.Type
        this.view.ce.DepComponent   = dependency.Component
        this.view.ce.DepVersion     = dependency.Version
        this.view.ce.Configuration2 = dependency.Configuration

        this.$forceUpdate()
      },
      // depTypeChanged handles changes of the type of the dependency information
      depTypeChanged: function() {
        // determine dependency
        dependency = this.model.Component.Dependencies[this.view.ce.Dependency]

        // update type
        dependency.Type = this.view.ce.DepType
      },
      // depComponentChanged handles changes of the component of the dependency information
      depComponentChanged: function() {
        // determine dependency
        dependency = this.model.Component.Dependencies[this.view.ce.Dependency]

        // update component
        dependency.Component = this.view.ce.DepComponent
      },
      // depVersionChanged handles changes of the version of the dependency information
      depVersionChanged: function() {
        // determine dependency
        dependency = this.model.Component.Dependencies[this.view.ce.Dependency]

        // update version
        dependency.Version = this.view.ce.DepVersion
      },
      // configuration2Changed handles changes of the dependency configuration information
      configuration2Changed: function() {
        // determine dependency
        dependency = this.model.Component.Dependencies[this.view.ce.Dependency]

        // update version
        dependency.Configuration = this.view.ce.Configuration2
      },
      // configuration2Focus expands the configuration editor of the dependency configuration
      configuration2Focus: function() {
        document.getElementById("ceConfiguration2").setAttribute("rows", 10)
      },
      // configuration2Blur schrinks the configuration editor of the dependency configuration
      configuration2Blur: function() {
        document.getElementById("ceConfiguration2").setAttribute("rows", 2)
      }
    },
    template: `
      <div class="modal-mask" v-if="model.Component != null">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Component Editor</h3>
            <div class="close" v-on:click="cancelDialog()"><i class="far fa-times-circle"></i></div>
          </div>

          <div class="modal-body">
            <table style="width: 100%">
              <col width="10*">
              <col width="990*">
              <col width="40px">
              <tr>
                <td><strong>Component:</strong></td>
                <td>
                  <input type="text"
                    v-model="view.ce.Component"
                    :disabled="!view.ce.New"
                    @change="componentChanged()"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Version:</td>
                <td>
                  <input type="text"
                    v-model="view.ce.Version"
                    @change="versionChanged()"
                    :disabled="view.ce.Component==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Element&nbsp;Configuration:</td>
                <td>
                  <textarea id="ceConfiguration1" rows=2
                    v-model="view.ce.Configuration1"
                    @change="configuration1Changed()"
                    @focus="configuration1Focus()"
                    @blur="configuration1Blur()"
                    :disabled="view.ce.Component==''"></textarea>
                </td>
              </tr>

              <tr>
                <td><strong>Dependency:</strong></td>
                <td>
                  <select
                    v-model="view.ce.Dependency"
                    @change="dependencyChanged()"
                    :disabled="view.ce.Component=='' || view.ce.Version==''">
                    <option disabled value="">please select</option>
                    <option v-for="d in dependencies()">{{d}}</option>
                  </select>
                </td>
                <td>
                  <div style="white-space: nowrap;">
                    <div class="icon" v-on:click="addDependency()"><i class="fas fa-plus-circle" ></i></div>
                    <div class="icon" v-on:click="delDependency()"><i class="fas fa-times-circle"></i></div>
                  </div>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency&nbsp;Type:</td>
                <td>
                  <select
                    v-model="view.ce.DepType"
                    @change="depTypeChanged()"
                    :disabled="view.ce.Component=='' || view.ce.Version==''">
                    <option disabled value="">please select</option>
                    <option>context</option>
                    <option>service</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency&nbsp;Component:</td>
                <td>
                  <select
                    v-model="view.ce.DepComponent"
                    @change="depComponentChanged()"
                    :disabled="view.ce.Component=='' || view.ce.Version==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in components()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency&nbsp;Version:</td>
                <td>
                  <select
                    v-model="view.ce.DepVersion"
                    @change="depVersionChanged()"
                    :disabled="view.ce.Component=='' || view.ce.Version==''">
                    <option disabled value="">please select</option>
                    <option v-for="v in versions()">{{v}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency&nbsp;Configuration:</td>
                <td>
                  <textarea id="ceConfiguration2"  rows=2
                    v-model="view.ce.Configuration2"
                    @change="configuration2Changed()"
                    @focus="configuration2Focus()"
                    @blur="configuration2Blur()"
                    :disabled="view.ce.Component=='' || view.ce.Version==''"></textarea>
                </td>
              </tr>

            </table>
          </div>

          <div class="modal-footer">
            &nbsp;
            <button class="modal-default-button"  v-if="!view.ce.New" v-on:click="duplicateComponent()">
              Duplicate Component <i class="fas fa-copy">
            </button>
            <button class="modal-default-button" v-if="!view.ce.New" v-on:click="deleteComponent()">
              Delete Component <i class="fas fa-times-circle">
            </button>
            <button class="modal-default-button" v-if="!view.ce.New" v-on:click="updateComponent()">
              Update Component <i class="fas fa-plus-circle">
            </button>
            <button class="modal-default-button" v-if="view.ce.New" v-on:click="createComponent()">
              Create Component <i class="fas fa-plus-circle">
            </button>
          </div>
        </div>
      </div>`
  }
)
