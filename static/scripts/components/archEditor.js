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
      // editConfiguration opens editor for editing a configuration
      editConfiguration: function(configuration_type, relationship) {
        // save the configuration type
        this.view.ae.ConfType = configuration_type

        switch( configuration_type) {
          case "element":
            this.view.ae.ConfTitle     = "Configuration of element '" + this.element.Element + "':"
            this.view.ae.Configuration = this.element.Configuration
            this.view.ae.Relationship  = ""
            break
          case "cluster":
            this.view.ae.ConfTitle     = "Configuration of cluster '" + this.view.ae.Cluster + "':"
            this.view.ae.Configuration = this.element.Clusters[this.view.ae.Cluster].Configuration
            this.view.ae.Relationship  = ""
            break
          case "relationship":
            this.view.ae.ConfTitle     = "Configuration of relationship: '" + relationship.Relationship + "':"
            this.view.ae.Configuration = relationship.Configuration
            this.view.ae.Relationship  = relationship.Relationship
            break
        }

        this.$forceUpdate()
      },
      // updateConfiguration updates the corresponding configuration
      updateConfiguration: function() {
        switch( this.view.ae.ConfType ) {
          case "element":
            this.model.ArchElement.Configuration = this.view.ae.Configuration
            break
          case "cluster":
            this.model.ArchElement.Clusters[this.view.ae.Cluster].Configuration = this.view.ae.Configuration
            break
          case "relationship":
            this.model.ArchElement.Clusters[this.view.ae.Cluster].Relationships[this.view.ae.Relationship] = this.view.ae.Configuration
            break
        }

        // close dialog
        this.discardConfiguration()
      },
      // discardConfiguration closes the configuration editor
      discardConfiguration: function() {
        this.view.ae.ConfType      = ""
        this.view.ae.ConfTitle     = ""
        this.view.ae.Relationship  = ""
        this.view.ae.Configuration = null

        this.$forceUpdate()
      },
      // addRelationship adds a new relationship to an element
      addRelationship: function() {
        // ask for name of new relationship and add it with initial values
        name = prompt("Name of the new relationship:")
        if (name != null && name != "" && name != "null") {
          // add relationship
          this.model.ArchElement.Clusters[view.ae.Cluster].Relationships[name] = {
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
        delete this.model.ArchElement.Clusters[view.ae.Cluster].Relationships[rel.Relationship]
        this.$forceUpdate()
      },
      // deleteElement removes an existing element from the architecture
      deleteElement: function() {
        Vue.delete(this.model.Architecture.Elements, this.model.ArchElement.Element)

        new ArchitectureGraph(this.model, this.view, this.view.domain, this.view.solution, this.view.version)

        this.model.ArchElement = null
      },
      // duplicateElement creates a copy of the element
      duplicateElement: function() {
        // ask for name of new element and create a copy
        name = prompt("Name of the new element:")
        if (name != null && name != "" && name != "null") {
          element = duplicateElement(this.model.ArchElement, name)

          // add element to architecture
          Vue.set(this.model.Architecture.Elements, name, element )

          // focus on new element
          this.model.ArchElement = this.model.Architecture.Elements[name]

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
            <td>&nbsp;Configuration:</td>
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
            <td>&nbsp;Configuration:</td>
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
                    <th>Configuration</th>
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

        <div class="configurationEditor" v-if="view.ae.Configuration != null">
          <div class="modal">
            <h3>{{view.ae.ConfTitle}}</h3>
            <textarea v-focus @keyup.esc="discardConfiguration()" v-model="view.ae.Configuration"></textarea>
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
