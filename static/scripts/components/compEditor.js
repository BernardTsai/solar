Vue.component(
  'compEditor',
  {
    props: ['model', 'view'],
    methods: {
      // addDependency adds a new dependency to a component
      addDependency: function() {
        // ask for name of new dependency and add it with initial values
        name = prompt("Name of the new dependency:")
        if (name != null && name != "" && name != "null") {
          // add dependency
          this.model.Component.Dependencies[name] = {
            Dependency:    name,
            Type:          "",
            Component:     "",
            Version:       "",
            Configuration: ""
          }
          this.$forceUpdate()
        }
      },
      // delDependency removes an existing dependency from a component
      delDependency: function(dep) {
        // remove dependency
        delete this.model.Component.Dependencies[dep.Dependency]
        this.$forceUpdate()
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
      // controllers lists the name and versions of available controllers
      controllers: function() {
        var result = []

        // loop over all available components
        for (var index in this.model.Controllers) {
          controller = this.model.Controllers[index]
          result.push( controller.Controller + " - " + controller.Version)
        }
        return result
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
      versions: function(comp) {
        var result = []

        // loop over all available components
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]

          // add new version to result
          if (component.Component == comp) {
            result.push(component.Version)
          }
        }
        return result
      },
      // componentChanged  handles changes of the component information
      componentChanged: function(dep) {
        this.$forceUpdate()
      },
      // editConfiguration opens editor for editing a configuration
      editConfiguration: function(dep) {
        if (dep != "") {
          this.view.ce.ConfTitle     = "Configuration template of dependency: " + dep.Dependency
          this.view.ce.Configuration = dep.Configuration
          this.view.ce.Dependency    = dep.Dependency
        } else {
          this.view.ce.ConfTitle     = "Configuration template of component: "
          this.view.ce.Configuration = this.model.Component.Configuration
          this.view.ce.Dependency    = ""
        }
        this.$forceUpdate()
      },
      // discardConfiguration closes the configuration editor
      discardConfiguration: function() {
        this.view.ce.Configuration = null
        this.view.ce.Dependency    = ""

        this.$forceUpdate()
      },
      // updateConfiguration updates the corresponding configuration
      updateConfiguration: function() {
        // update configuration
        depName = this.view.ce.Dependency
        if (depName != "") {
          dep = this.model.Component.Dependencies[depName]
          dep.Configuration = this.view.ce.Configuration
        } else {
          this.model.Component.Configuration = this.view.ce.Configuration
        }
        // close dialog
        this.view.ce.Configuration = null
        this.view.ce.Dependency    = ""

        this.$forceUpdate()
      }
    },
    template: `
      <div class="compEditor" v-if="model.Component">
        <div class="header">
          <h3>Component: {{model.Component.Component}} ({{model.Component.Version}})</h3>
          <button class="modal-default-button"  v-if="!view.ce.New" v-on:click="duplicateComponent()">
            Duplicate <i class="fas fa-copy">
          </button>
          <button class="modal-default-button" v-if="!view.ce.New" v-on:click="deleteComponent()">
            Delete <i class="fas fa-times-circle">
          </button>
          <button class="modal-default-button" v-if="!view.ce.New" v-on:click="updateComponent()">
            Update <i class="fas fa-plus-circle">
          </button>
          <button class="modal-default-button" v-if="view.ce.New" v-on:click="createComponent()">
            Create <i class="fas fa-plus-circle">
          </button>
        </div>

        <table style="width: 100%">
          <col width="10*">
          <col width="990*">
          <tr>
            <td><strong>Component:</strong></td>
            <td>
              <input type="text"
                v-model="model.Component.Component">
              </input>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Version:</td>
            <td>
              <input type="text"
                v-model="model.Component.Version">
              </input>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Controller:</td>
            <td>
              <select
                v-model="model.Component.Controller">
                <option disabled value="">please select</option>
                <option v-for="c in controllers()">{{c}}</option>
              </select>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Conf.&nbsp;Template:</td>
            <td>
              <textarea id="configuration" rows=10
                @focus="editConfiguration('')"
                v-model="model.Component.Configuration">
              </textarea>
            </td>
          </tr>
          <tr>
            <td><strong>Dependencies:</strong></td>
            <td>

              <table class="dependencies">
                <thead>
                  <tr>
                    <th>Dependency</th>
                    <th>Type</th>
                    <th>Component</th>
                    <th>Version</th>
                    <th>Conf.&nbsp;Template</th>
                    <th class="center" @click="addDependency"><i class="fas fa-plus-circle"></i></th>
                </thead>
                <tbody>
                  <tr v-for="dep in model.Component.Dependencies">
                    <td>
                      <input type="text"
                        v-model="dep.Dependency">
                      </input>
                    </td>
                    <td>
                      <select
                        v-model="dep.Type">
                        <option disabled value="">please select</option>
                        <option>context</option>
                        <option>service</option>
                      </select>
                    </td>
                    <td>
                      <select
                        @change="componentChanged(dep)"
                        v-model="dep.Component">
                        <option disabled value="">please select</option>
                        <option v-for="c in components()">{{c}}</option>
                      </select>
                    </td>
                    <td>
                      <select
                        v-model="dep.Version">>
                        <option disabled value="">please select</option>
                        <option v-for="v in versions(dep.Component)">{{v}}</option>
                      </select>
                    </td>
                    <td @click="editConfiguration(dep)">
                      <i class="fas fa-edit"></i>
                    </td>
                    <td @click="delDependency(dep)">
                      <i class="fas fa-minus-circle"></i>
                    </td>
                  </tr>
                </tbody>
              </table>

            </td>
          </tr>
        </table>

        <div class="configurationEditor" v-if="view.ce.Configuration != null">
          <div class="modal">
            <h3>{{view.ce.ConfTitle}}</h3>
            <textarea v-focus @keyup.esc="discardConfiguration()" v-model="view.ce.Configuration"></textarea>
            <button class="modal-default-button" @click="updateConfiguration()">
              OK <i class="fas fa-check-circle">
            </button>
            <button class="modal-default-button" @click="discardConfiguration()">
              Cancel <i class="fas fa-times-circle">
            </button>

          </div>
        </div>

      </div>`
  }
)
