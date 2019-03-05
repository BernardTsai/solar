Vue.component(
  'architectureElementWizard',
  {
    props: ['model', 'view','name','component','element'],
    computed: {
      components: function() {
        var map = {};
        for (var index in model.Catalog) {
          component = model.Catalog[index]
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
      // clusters
      clusters: function() {
        // check if the name has been defined
        if (!this.element.component) {
          return []
        }
        // match versions
        var map = {};
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
          name      = component.Component
          version   = component.Version

          if (name == this.element.component) {
            if (!map[version]) {
              console.log(version)
              map[version] = version
            }
          }
        }
        var result = []
        for (var version in map) {
          result.push(version)
        }
        return result
      },
      // relationships
      relationships: function() {
        result = []
        return result
      },
      // dependencies
      dependencies: function() {
        // check if the name and version has been defined
        if (!this.element.component || !this.element.cluster) {
          return []
        }
        // match relationships
        var result = []
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
          name      = component.Component
          version   = component.Version

          if (name == this.element.component && version == this.element.version) {
            for (var rel in component.Dependencies) {
              result.push(rel.Dependency)
            }
          }
        }
        return result
      },
      // componentsChanged updates the computation
      componentChanged: function() {
        // add possible clusters to the element
        for (index in model.Catalog) {
          comp = model.Catalog[index]

          if (comp.Component == this.element.component) {
            cluster = {
              Version:       comp.Version,
              State:         "active",
              Min:           1,
              Max:           1,
              Size:          1,
              Configuration: "",
              Relationships: {}
            }

            // add relationships to the clusters
            for (index2 in comp.Dependencies) {
              dependency = comp.Dependencies[index2]

              relationship = {
                Relationship:  dependency.Dependency,
                Dependency:    dependency.Dependency,
                Type:          dependency.Type,
                Element:       "",
                Version:       "",
                Configuration: ""
              }

              // add relationship to cluster
              cluster.Relationships[relationship.Relationship] = relationship
            }

            // add cluster to element
            view.element.Clusters[cluster.Version] = cluster
          }
        }
        this.$forceUpdate()
      },
      // cancelDialog closes the wizard dialog
      cancelDialog: function() {
        this.view.element    = null
        this.view.newElement = false
      },
      // createElement constructs a new element for the architecture
      createElement: function() {
        // construct base element
        element = {
          Element:       this.name,
          Component:     this.component,
          Configuration: "",
          Clusters:      {}
        }

        // add possible clusters to the element
        for (index in model.Catalog) {
          comp = model.Catalog[index]

          if (comp.Component == this.component) {
            cluster = {
              Version:       comp.Version,
              State:         "active",
              Min:           1,
              Max:           1,
              Size:          1,
              Configuration: "",
              Relationships: {}
            }

            // add relationships to the clusters
            for (index2 in comp.Dependencies) {
              dependency = comp.Dependencies[index2]

              relationship = {
                Relationship:  dependency.Dependency,
                Dependency:    dependency.Dependency,
                Type:          dependency.Type,
                Element:       "",
                Version:       "",
                Configuration: ""
              }

              // add relationship to cluster
              cluster.Relationships[relationship.Relationship] = relationship
            }

            // add cluster to element
            element.Clusters[cluster.Version] = cluster
          }
        }

        // open editor for configuring the element
        this.view.Element    = element
        this.view.newElement = true
        this.$forceUpdate()
      }
    },
    template: `
      <div class="modal-mask" v-if="view.element != null">
        <div class="modal-wrapper">

          <div class="modal-container">

            <div class="modal-header">
              <h3>Architecture Element Wizard</h3>
              <div class="close" v-on:click="cancelDialog()"><i class="far fa-times-circle"></i></div>
            </div>

            <div class="modal-body">
              <table style="width: 100%">
                <col width="10*">
                <col width="990*">
                <col width="40px">
                <tr>
                  <td><strong>Element:</strong></td>
                  <td>
                    <input class="long" type="text" v-model="element.name"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Component:</td>
                  <td>
                    <select v-model="element.component" v-on:change="componentChanged()">
                      <option disabled value="">component</option>
                      <option v-for="c in components">{{c}}</option>
                    </select>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Element&nbsp;Configuration:</td>
                  <td>
                    <input class="long" type="text" v-model="element.configuration1"></input>
                  </td>
                </tr>

                <tr>
                  <td><strong>Cluster:</strong></td>
                  <td>
                    <select v-model="element.cluster">
                      <option disabled value="">cluster</option>
                      <option v-for="c in clusters()">{{c}}</option>
                    </select>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;State:</td>
                  <td>
                    <input class="long" type="text" v-model="element.state"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Min:</td>
                  <td>
                    <input class="long" type="text" v-model="element.min"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Max:</td>
                  <td>
                    <input class="long" type="text" v-model="element.max"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Size:</td>
                  <td>
                    <input class="long" type="text" v-model="element.size"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Cluster&nbsp;Configuration:</td>
                  <td>
                    <input class="long" type="text" v-model="element.configuration2"></input>
                  </td>
                </tr>

                <tr>
                  <td><strong>Relationship:</strong></td>
                  <td>
                    <select v-model="element.relationship">
                      <option disabled value="">relationship</option>
                      <option v-for="c in relationships()">{{c}}</option>
                    </select>
                  </td>
                  <td>
                    <div style="white-space: nowrap;">
                      <i class="fas fa-plus-circle"></i>
                      <i class="fas fa-minus-circle"></i>
                    </div>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Dependency:</td>
                  <td>
                    <select v-model="element.dependency">
                      <option disabled value="">dependency</option>
                      <option v-for="c in dependencies()">{{c}}</option>
                    </select>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Type:</td>
                  <td>
                    <select v-model="element.relType">
                      <option disabled value="">please select</option>
                      <option>context</option>
                      <option>service</option>
                    </select>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Element:</td>
                  <td>
                    <input class="long" type="text" v-model="element.relElement"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Version:</td>
                  <td>
                    <input class="long" type="text" v-model="element.relVersion"></input>
                  </td>
                </tr>
                <tr>
                  <td>&nbsp;Relationship&nbsp;Configuration:</td>
                  <td>
                    <input class="long" type="text" v-model="element.relConfiguration"></input>
                  </td>
                </tr>

              </table>
            </div>

            <div class="modal-footer">
              &nbsp;
              <button class="modal-default-button" v-on:click="createElement()">
                Create Element <i class="fas fa-check">
              </button>
            </div>
          </div>
        </div>
      </div>`
  }
)
