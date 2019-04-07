Vue.component(
  'archEditor',
  {
    props: ['model', 'view','element'],
    methods: {
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
      // dependencies list the possible dependencies the element can have
      dependencies: function(component, version) {
        var result = []

        // find matching component
        for (var c of this.model.Catalog) {
          // add all dependencies
          if (c.Component == component && c.Version == version) {
            result = Object.keys(c.Dependencies)
            break
          }
        }

        // return sorted result
        result.sort()

        return result
      },
      // relDependencyChanged handles the change of a relationship dependency
      relDependencyChanged(component, rel) {
        // check if a dependency has been selected
        if (rel.Dependency == "") {
          rel.Version = ""
        } else {
          // add version to relationship
          rel.Version = ""
          for (var c of this.model.Catalog) {
            if (c.Component == component) {
              for (var d in c.Dependencies) {
                if (d == rel.Dependency) {
                  rel.Type    = c.Dependencies[d].Type
                  rel.Version = c.Dependencies[d].Version
                  break
                }
              }
              // dependency has already been found
              if (rel.Version != "") {
                break
              }
            }
          }
        }

        // update GUI
        this.$forceUpdate()
      },
      // elements lists all available elements
      elements: function(component, version, dependency) {
        // check parameters
        if (component == "" || version == "" || dependency == "") {
          return []
        }

        // find dependency
        dep = null
        for (var c of this.model.Catalog) {
          if (c.Component == component && c.Version == version) {
            dep = c.Dependencies[dependency]
            break
          }
        }

        // check if dependency was found
        if (!dep) {
          return []
        }

        // list eleents
        result = []
        Object.values(this.model.Architecture.Elements).forEach((element) => {
          if (element.Component == dep.Component) {
            result.push(element.Element)
          }
        })

        // return the sorted result
        result.sort()

        return result
      },
      // alignParameters initialises the configuration parameters
      alignParameters: function(configuration_type, template, configuration) {
        dict = jsyaml.safeLoad(configuration)
        if (!dict) {
          dict = {}
        }
        dict["domain"]       = this.view.domain
        dict["solution"]     = this.view.solution
        dict["version"]      = this.view.version
        dict["element"]      = this.element.Element
        dict["cluster"]      = this.view.ae.Cluster
        dict["relationship"] = this.view.ae.Relationship

        defDict = {"domain": 0, "solution":0, "version":0, "element":0, "cluster": 0, "relationship": 0}
        config = "# Configuration parameters\n"
        names = []
        length = 12
        matches = template.match(/{{[^}]*}}/g)
        if (matches) {
          for (var name of matches) {
            name = name.replace("{{","")
            name = name.replace("}}","")
            name = name.trim()
            if (!(name in defDict)) {
              defDict[name] = null
              names.push(name)
              length = Math.max(length, name.length)
              if (!(name in dict)) {
                dict[name] = null
              }
            }
          }
        }
        if (configuration_type == 'relationship') {
          names.unshift("relationship")
        }
        if (configuration_type != 'element') {
          names.unshift("cluster")
        }
        names.sort()
        names.unshift("element")
        names.unshift("version")
        names.unshift("solution")
        names.unshift("domain")

        // construct parameters
        config = "# " + this.view.ae.ConfTitle + "\n"
        for (var name of names) {
          key = "'" + name + "': " + " ".repeat(length)
          key = key.substr(0, length + 4)
          if (dict[name]) {
            config += key + "'" + dict[name] + "'\n"
          } else {
            config += key + "'<enter parameter here>'\n"
          }
        }

        return config
      },
      // showTemplate displays the template
      showTemplate: function() {
        this.view.ae.Display = "template"
      },
      // showParameters displays the parameters
      showParameters: function() {
        this.view.ae.Display = "parameters"
      },
      // showConfiguration renders the component template with the parameters
      showConfiguration: function() {
        configuration = this.view.ae.Template
        parameters    = jsyaml.safeLoad(this.view.ae.Parameters)

        repeat = true
        while (repeat) {
          repeat = false
          matches = configuration.match(/{{[^}]*}}/g)
          if (matches) {
            for (var name of matches) {
              key = name
              key = key.replace("{{","")
              key = key.replace("}}","")
              key = key.trim()
              if (key in parameters) {
                configuration = configuration.replace(new RegExp(name, 'g'), parameters[key])
                repeat = true
              }
            }
          }
        }

        this.view.ae.Configuration = configuration
        this.view.ae.Display       = "configuration"
      },
      // editConfiguration opens editor for editing a configuration
      editConfiguration: function(configuration_type, relationship) {
        // save the configuration type
        this.view.ae.ConfType = configuration_type

        switch( configuration_type) {
          case "element":
            this.view.ae.ConfTitle     = "Configuration for element '" + this.element.Element + "':"
            this.view.ae.Display       = "parameters"
            this.view.ae.Relationship  = ""
            // find templates of all matching component versions
            templates = ""
            for (var c of this.model.Catalog) {
              // add all dependencies
              if (c.Component == this.element.Component) {
                templates += "-".repeat(80) + "\n"
                templates += "Configuration template of component version: " + c.Version + "\n"
                templates += "-".repeat(80) + "\n\n"
                templates += c.Configuration + "\n\n"
              }
            }
            this.view.ae.Template      = templates
            this.view.ae.Parameters    = this.alignParameters(configuration_type, templates, this.element.Configuration)
            this.view.ae.Configuration = ""
            break
          case "cluster":
            this.view.ae.ConfTitle     = "Configuration for cluster '" + this.view.ae.Cluster + "':"
            this.view.ae.Display       = "parameters"
            this.view.ae.Relationship  = ""
            // find templates of the matching component version
            template = ""
            for (var c of this.model.Catalog) {
              // add all dependencies
              if (c.Component == this.element.Component && c.Version == this.view.ae.Cluster) {
                template += "-".repeat(80) + "\n"
                template += "Configuration template of component version: " + c.Version + "\n"
                template += "-".repeat(80) + "\n\n"
                template += c.Configuration + "\n\n"
                break
              }
            }
            this.view.ae.Template      = template
            this.view.ae.Parameters    = this.alignParameters(configuration_type, template, this.element.Clusters[this.view.ae.Cluster].Configuration)
            this.view.ae.Configuration = ""
            break
          case "relationship":
            this.view.ae.ConfTitle     = "Configuration for relationship: '" + relationship.Relationship + "':"
            this.view.ae.Display       = "parameters"
            this.view.ae.Relationship  = relationship.Relationship
            // find templates of the matching dependency
            template = ""
            for (var c of this.model.Catalog) {
              if (c.Component == this.element.Component && c.Version == this.view.ae.Cluster) {
                dep = c.Dependencies[relationship.Dependency]
                template += "-".repeat(80) + "\n"
                template += "Configuration template of dependency : " + dep.Dependency + "\n"
                template += "-".repeat(80) + "\n\n"
                template += dep.Configuration + "\n\n"
                break
              }
            }
            this.view.ae.Template      = template
            this.view.ae.Parameters    = this.alignParameters(configuration_type, template, this.element.Clusters[this.view.ae.Cluster].Relationships[this.view.ae.Relationship].Configuration)
            this.view.ae.Configuration = ""
            break
        }

        this.$forceUpdate()
      },
      // updateConfiguration updates the corresponding configuration
      updateConfiguration: function() {
        switch( this.view.ae.ConfType ) {
          case "element":
            this.element.Configuration = this.view.ae.Parameters
            break
          case "cluster":
            this.element.Clusters[this.view.ae.Cluster].Configuration = this.view.ae.Parameters
            break
          case "relationship":
            this.element.Clusters[this.view.ae.Cluster].Relationships[this.view.ae.Relationship].Configuration = this.view.ae.Parameters
            break
        }

        // close dialog
        this.discardConfiguration()
      },
      // discardConfiguration closes the configuration editor
      discardConfiguration: function() {
        this.view.ae.Display       = ""
        this.view.ae.Template      = ""
        this.view.ae.Parameters    = ""
        this.view.ae.Configuration = ""
        this.view.ae.ConfType      = ""
        this.view.ae.ConfTitle     = ""
        this.view.ae.Relationship  = ""

        this.$forceUpdate()
      },
      // addRelationship adds a new relationship to an element
      addRelationship: function() {
        // ask for name of new relationship and add it with initial values
        name = prompt("Name of the new relationship:")
        if (name != null && name != "" && name != "null") {
          // add relationship
          this.element.Clusters[view.ae.Cluster].Relationships[name] = {
            Relationship:  name,
            Dependency:    "",
            Element:       "",
            Version:       "",
            Configuration: ""
          }
          this.$forceUpdate()
        }
      },
      // delRelationship removes an existing relationship from an element
      delRelationship: function(rel) {
        // remove dependency
        delete this.element.Clusters[view.ae.Cluster].Relationships[rel.Relationship]
        this.$forceUpdate()
      },
      // deleteElement removes an existing element from the architecture
      deleteElement: function() {
        Vue.delete(this.model.Architecture.Elements, this.element.Element)

        new ArchitectureGraph(this.model, this.view, this.view.domain, this.view.solution, this.view.version)

        this.element = null
      },
      // duplicateElement creates a copy of the element
      duplicateElement: function() {
        // ask for name of new element and create a copy
        name = prompt("Name of the new element:")
        if (name != null && name != "" && name != "null") {
          element = duplicateElement(this.element, name)

          // add element to architecture
          Vue.set(this.model.Architecture.Elements, name, element )

          // focus on new element
          this.element = this.model.Architecture.Elements[name]

          // force update
          this.$forceUpdate()
        }
      }
    },
    template: `
      <div class="archEditor">
        <div class="header">
          <h3>Element: {{element.Element}} ({{element.Component}})</h3>

          <button class="modal-default-button" v-if="!view.ce.New" v-on:click="duplicateElement()">
            Duplicate <i class="fas fa-copy">
          </button>

          <button class="modal-default-button" v-if="!view.ce.New" v-on:click="deleteElement()">
            Delete <i class="fas fa-times-circle">
          </button>
        </div>

        <table style="width: 100%">
          <col width="10*">
          <col width="990*">
          <tr>
            <td><strong>Element:</strong></td>
            <td>
              <input type="text" readonly v-model="element.Element"/>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Component:</td>
            <td>
              <input type="text" readonly v-model="element.Component"/>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Conf.&nbsp;Parameters:</td>
            <td>
              <textarea id="configuration" rows=5
                @click="editConfiguration('element', '')"
                v-model="element.Configuration">
              </textarea>
            </td>
          </tr>

          <tr>
            <td><strong>Clusters:</strong></td>
            <td>
              <select
                v-model="view.ae.Cluster">
                <option value="">please select</option>
                <option v-for="cluster in element.Clusters">{{cluster.Version}}</option>
              </select>
            </td>
          </tr>
          <tr v-if="view.ae.Cluster != ''">
            <td>&nbsp;Target State:</td>
            <td>
              <select
                v-model="element.Clusters[view.ae.Cluster].State">
                <option disabled value="">please select</option>
                <option>initial</option>
                <option>inactive</option>
                <option>active</option>
              </select>
            </td>
          </tr>
          <tr v-if="view.ae.Cluster != ''">
            <td>&nbsp;Minimum Size:</td>
            <td>
              <input type="number" v-model="element.Clusters[view.ae.Cluster].Min"/>
            </td>
          </tr>
          <tr v-if="view.ae.Cluster != ''">
            <td>&nbsp;Maximum Size:</td>
            <td>
              <input type="number" v-model="element.Clusters[view.ae.Cluster].Max"/>
            </td>
          </tr>
          <tr v-if="view.ae.Cluster != ''">
            <td>&nbsp;Current Size:</td>
            <td>
              <input type="number" v-model="element.Clusters[view.ae.Cluster].Size"/>
            </td>
          </tr>
          <tr v-if="view.ae.Cluster != ''">
            <td>&nbsp;Conf.&nbsp;Parameters:</td>
            <td>
              <textarea id="configuration2" rows=5 readonly
                @click="editConfiguration('cluster', '')"
                v-model="element.Clusters[view.ae.Cluster].Configuration">
              </textarea>
            </td>
          </tr>

          <tr v-if="view.ae.Cluster != ''">
            <td><strong>Relationships:</strong></td>
            <td>

              <table class="relationships">
                <thead>
                  <tr>
                    <th>Relationship</th>
                    <th>Dependency</th>
                    <th>Element</th>
                    <th>Conf.&nbsp;Parameters</th>
                    <th class="center" @click="addRelationship"><i class="fas fa-plus-circle"></i></th>
                </thead>
                <tbody>
                  <tr v-for="rel in element.Clusters[view.ae.Cluster].Relationships">
                    <td>
                      <input type="text" v-model="rel.Relationship"/>
                    </td>
                    <td>
                      <select
                        @change="relDependencyChanged(element.Component, rel)"
                        @change="$forceUpdate()"
                        v-model="rel.Dependency">
                        <option disabled value="">please select</option>
                        <option v-for="dep in dependencies(element.Component, view.ae.Cluster)">{{dep}}</option>
                      </select>
                    </td>
                    <td>
                      <select
                        v-model="rel.Element">
                        <option disabled value="">please select</option>
                        <option v-for="element in elements(element.Component, view.ae.Cluster, rel.Dependency)">{{element}}</option>
                      </select>
                    </td>
                    <td @click="editConfiguration('relationship', rel)">
                      <i class="fas fa-edit"></i>
                    </td>
                    <td @click="delRelationship(rel)">
                      <i class="fas fa-minus-circle"></i>
                    </td>
                  </tr>
                </tbody>
              </table>

            </td>
          </tr>
        </table>

        <div class="configurationEditor" v-if="view.ae.Display != ''">
          <div class="modal">
            <h3>{{view.ae.ConfTitle}}</h3>
            <textarea v-focus readonly @keyup.esc="discardConfiguration()" v-if="view.ae.Display=='configuration'" v-model="view.ae.Configuration"></textarea>
            <textarea v-focus          @keyup.esc="discardConfiguration()" v-if="view.ae.Display=='parameters'"    v-model="view.ae.Parameters"></textarea>
            <textarea v-focus readonly @keyup.esc="discardConfiguration()" v-if="view.ae.Display=='template'"      v-model="view.ae.Template"></textarea>
            <button class="modal-default-button" @click="updateConfiguration()">
              OK <i class="fas fa-check-circle">
            </button>
            <button class="modal-default-button" @click="discardConfiguration()">
              Cancel <i class="fas fa-times-circle">
            </button>
            <button class="modal-default-button" @click="showConfiguration()">
              Configuration <i class="fas fa-file-alt">
            </button>
            <button class="modal-default-button" @click="showParameters()">
              Parameters <i class="fas fa-align-justify">
            </button>
            <button class="modal-default-button" @click="showTemplate()">
              Template <i class="fas fa-file-code">
            </button>

          </div>
        </div>

      </div>`
  }
)
