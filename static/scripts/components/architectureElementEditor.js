Vue.component(
  'architectureElementEditor',
  {
    props: ['model', 'view'],
    methods: {
      addRelationship: function() {
        // check if a cluster has been selected
        if (this.view.ae.Cluster=="") {
          return
        }

        // ask for name of new relationship and add it with initial values
        name = prompt("Name of the new relationship:")
        if (name && name != "") {
          // add relationship
          this.model.ArchElement.Clusters[this.view.ae.Cluster].Relationships[name] = {
            Relationship:   name,
            Dependency:     "",
            RelElement:     "",
            Configuration3: ""
          }

          // update dialog
          this.view.ae.Relationship   = name
          this.view.ae.Dependency     = ""
          this.view.ae.RelElement     = ""
          this.view.ae.Configuration3 = ""
        }
      },
      delRelationship: function() {
        // check if relationship has been selected
        if (this.view.ae.Relationship=="") {
          return
        }

        // remove relationship
        delete this.model.ArchElement.Clusters[this.view.ae.Cluster].Relationships[this.view.ae.Relationship]

        // update the fields
        this.view.ae.Relationship=   ""
        this.view.ae.Dependency=     ""
        this.view.ae.RelElement=     ""
        this.view.ae.Configuration3= ""
      },
      configuration1Focus: function() {
        document.getElementById("aeConfiguration1").setAttribute("rows", 10)
      },
      configuration1Blur: function() {
        document.getElementById("aeConfiguration1").setAttribute("rows", 2)
      },
      configuration2Focus: function() {
        document.getElementById("aeConfiguration2").setAttribute("rows", 10)
      },
      configuration2Blur: function() {
        document.getElementById("aeConfiguration2").setAttribute("rows", 2)
      },
      configuration3Focus: function() {
        document.getElementById("aeConfiguration3").setAttribute("rows", 10)
      },
      configuration3Blur: function() {
        document.getElementById("aeConfiguration3").setAttribute("rows", 2)
      },
      // components creates a list of component names from the catalog
      components: function() {
        var map = {};
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
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
      },
      // clusters creates a list of available cluster versions
      clusters: function() {
        var result = []
        for (var index in this.model.ArchElement.Clusters) {
          version = this.model.ArchElement.Clusters[index].Version
          result.push(version)
        }
        return result
      },
      // relationships
      relationships: function() {
        // check if the name and version and relationship have been defined
        if (this.view.ae.Component == "" || this.view.ae.Cluster == "") {
          return []
        }
        // match relationships
        cluster = this.model.ArchElement.Clusters[this.view.ae.Cluster]
        var result = []
        for (var index in cluster.Relationships) {
          relationship = cluster.Relationships[index]
          result.push(relationship.Relationship)
        }
        return result
      },
      // dependencies
      dependencies: function() {
        // check if the name and version and relationship have been defined
        if (this.view.ae.Component == "" || this.view.ae.Cluster == "" ||  this.view.ae.Relationship == "") {
          return []
        }
        var map = {};

        // loop over all components of the catalog
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
          name      = component.Component
          version   = component.Version

          // check if name and version match
          if (name == this.view.ae.Component && version == this.view.ae.Cluster) {
            // loop over all dependencies
            for (var index2 in component.Dependencies) {
              dependency = component.Dependencies[index2]

              map[dependency.Dependency] = dependency.Dependency
            }
            break
          }
        }
        var result = []
        for (var dependency in map) {
          result.push(dependency)
        }
        return result
      },
      // elements
      elements: function() {
        // check if the name and version and relationship have been defined
        if (this.view.ae.Component == "" || this.view.ae.Cluster == "" ||  this.view.ae.Relationship == "" || this.view.ae.Dependency == "") {
          return []
        }

        // determine dependency and corresponding component type and version of dependency
        for (var index in this.model.Catalog) {
          component = this.model.Catalog[index]
          name      = component.Component
          version   = component.Version

          // check if name and version match
          if (name == this.view.ae.Component && version == this.view.ae.Cluster) {
            dep = component.Dependencies[this.view.ae.Dependency]

            depComponent = dep.Component
            depVersion   = dep.Version

            break
          }
        }

        // compile all matching elements
        result = []
        for (var index in this.model.Architecture.Elements) {
          element = this.model.Architecture.Elements[index]

          if (element.Component == depComponent && element.Clusters[depVersion]) {
            result.push(element.Element)
          }
        }
        return result
      },
      // elementChanged  handles changes of the element information
      elementChanged: function() {
        this.model.ArchElement.Element = this.view.ae.Element
      },
      // componentChanged  handles changes of the component information
      componentChanged: function() {
        this.model.ArchElement.Component = this.view.ae.Component

        // add possible clusters to the element
        this.model.ArchElement.Clusters = {}
        for (index in this.model.Catalog) {
          comp = this.model.Catalog[index]

          if (comp.Component == this.view.ae.Component) {
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
                Relationship:  dependency.Dependency + " - Relationship",
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
            this.model.ArchElement.Clusters[cluster.Version] = cluster
          }
        }

        // update the fields
        this.view.ae.Cluster=        ""
        this.view.ae.State=          ""
        this.view.ae.Min=            ""
        this.view.ae.Max=            ""
        this.view.ae.Size=           ""
        this.view.ae.Configuration2= ""
        this.view.ae.Relationship=   ""
        this.view.ae.Dependency=     ""
        this.view.ae.RelElement=     ""
        this.view.ae.Configuration3= ""

        this.$forceUpdate()
      },
      // configuration1Changed  handles changes of the element configuration information
      configuration1Changed: function() {
        this.model.ArchElement.Configuration = this.view.ae.Configuration1
      },
      // clusterChanged handles changes of the cluster information
      clusterChanged: function() {
        // determine cluster
        cluster = this.model.ArchElement.Clusters[this.view.ae.Cluster]

        // update the fields
        this.view.ae.State          = (cluster.State != "" ? cluster.State : "active")
        this.view.ae.Min            = (cluster.Min   != "" ? cluster.Min   : 1)
        this.view.ae.Max            = (cluster.Max   != "" ? cluster.Max   : 1)
        this.view.ae.Size           = (cluster.Size  != "" ? cluster.Size  : 1)
        this.view.ae.Configuration2 = cluster.Configuration
        this.view.ae.Relationship   = ""
        this.view.ae.Dependency     = ""
        this.view.ae.RelElement     = ""
        this.view.ae.Configuration3 = ""

        this.$forceUpdate()
      },
      // stateChanged  handles changes of the state information
      stateChanged: function() {
        this.model.ArchElement.Clusters[this.view.ae.Cluster].Status = this.view.ae.Status
      },
      // minChanged  handles changes of the min information
      minChanged: function() {
        this.model.ArchElement.Clusters[this.view.ae.Cluster].Min = this.view.ae.Min
      },
      // maxChanged  handles changes of the max information
      maxChanged: function() {
        this.model.ArchElement.Clusters[this.view.ae.Cluster].Max = this.view.ae.Max
      },
      // sizeChanged  handles changes of the size information
      sizeChanged: function() {
        this.model.ArchElement.Clusters[this.view.ae.Cluster].Size = this.view.ae.Size
      },
      // configuration2Changed  handles changes of the cluster configuration information
      configuration2Changed: function() {
        this.model.ArchElement.Clusters[this.view.ae.Cluster].Configuration = this.view.ae.Configuration2
      },
      // relationshipChanged handles changes of the relationship information
      relationshipChanged: function() {
        // determine cluster and relationship
        cluster = this.model.ArchElement.Clusters[this.view.ae.Cluster]
        relationship = cluster.Relationships[this.view.ae.Relationship]

        // update the fields
        this.view.ae.Dependency=     relationship.Dependency
        this.view.ae.RelElement=     relationship.Element
        this.view.ae.Configuration3= relationship.Configuration

        this.$forceUpdate()
      },
      // dependencyChanged handles changes of the relationship dependency information
      dependencyChanged: function() {
        // determine cluster and relationship
        cluster = this.model.ArchElement.Clusters[this.view.ae.Cluster]
        relationship = cluster.Relationships[this.view.ae.Relationship]

        // update the field
        relationship.Dependency = this.view.ae.Dependency

        // reset element and configuration
        this.view.ae.RelElement     = ""
        this.view.ae.Configuration3 = ""
      },
      // relElementChanged handles changes of the relationship element information
      relElementChanged: function() {
        // determine cluster and relationship and update element
        cluster              = this.model.ArchElement.Clusters[this.view.ae.Cluster]
        relationship         = cluster.Relationships[this.view.ae.Relationship]
        relationship.Element = this.view.ae.RelElement
      },
      // configuration3Changed handles changes of the relationship configuration information
      configuration3Changed: function() {
        // determine cluster and relationship and update configuration
        cluster                    = this.model.ArchElement.Clusters[this.view.ae.Cluster]
        relationship               = cluster.Relationships[this.view.ae.Relationship]
        relationship.Configuration = this.view.ae.Configuration3
      },

      // cancelDialog closes the wizard dialog
      cancelDialog: function() {
        this.model.ArchElement = null
        this.view.ae.New       = false
      },
      // createElement adds element to the architecture
      createElement: function() {
        Vue.set(this.model.Architecture.Elements, this.view.ae.Element, this.model.ArchElement)
        this.model.ArchElement = null
      },
      // updateElement updates an existing element of the architecture
      updateElement: function() {
        Vue.delete(this.model.Architecture.Elements, this.view.ae.Element)
        Vue.set(this.model.Architecture.Elements, this.view.ae.Element, this.model.ArchElement)
        this.model.ArchElement = null
      },
      // deleteElement removes an existing element from the architecture
      deleteElement: function() {
        Vue.delete(this.model.Architecture.Elements, this.view.ae.Element)
        this.model.ArchElement = null
      }
    },
    template: `
      <div class="modal-mask" v-if="model.ArchElement != null">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Architecture Element Editor</h3>
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
                  <input type="text" v-model="view.ae.Element" :disabled="!view.ae.New" @change="elementChanged()"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Component:</td>
                <td>
                  <select v-model="view.ae.Component" v-on:change="componentChanged()"
                    :disabled="view.ae.Element==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in components()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Element&nbsp;Configuration:</td>
                <td>
                  <textarea id="aeConfiguration1" rows=2
                    v-model="view.ae.Configuration1"
                    @change="configuration1Changed()"
                    @focus="configuration1Focus()"
                    @blur="configuration1Blur()"
                    :disabled="view.ae.Element==''"></textarea>
                </td>
              </tr>

              <tr>
                <td><strong>Cluster:</strong></td>
                <td>
                  <select v-model="view.ae.Cluster" v-on:change="clusterChanged()" :disabled="view.ae.Component==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in clusters()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;State:</td>
                <td>
                  <select v-model="view.ae.State" v-on:change="stateChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''">
                    <option disabled value="">please select</option>
                    <option>initial</option>
                    <option>inactive</option>
                    <option>active</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Min:</td>
                <td>
                  <input class="long" type="number" min="0" v-model="view.ae.Min" v-on:change="minChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Max:</td>
                <td>
                  <input class="long" type="number" min="0" v-model="view.ae.Max" v-on:change="maxChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Size:</td>
                <td>
                  <input class="long" type="number" min="0" v-model="view.ae.Size" v-on:change="sizeChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Cluster&nbsp;Configuration:</td>
                <td>
                  <textarea id="aeConfiguration2" rows=2
                    v-model="view.ae.Configuration2"
                    @change="configuration2Changed()"
                    @focus="configuration2Focus()"
                    @blur="configuration2Blur()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''"></textarea>
                </td>
              </tr>

              <tr>
                <td><strong>Relationship:</strong></td>
                <td>
                  <select v-model="view.ae.Relationship" v-on:change="relationshipChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in relationships()">{{c}}</option>
                  </select>
                </td>
                <td>
                  <div style="white-space: nowrap;">
                    <div class="icon" v-on:click="addRelationship()"><i class="fas fa-plus-circle" ></i></div>
                    <div class="icon" v-on:click="delRelationship()"><i class="fas fa-times-circle"></i></div>
                  </div>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency:</td>
                <td>
                  <select v-model="view.ae.Dependency"  v-on:change="dependencyChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster=='' || view.ae.Relationship==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in dependencies()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Element:</td>
                <td>
                  <select v-model="view.ae.RelElement"  v-on:change="relElementChanged()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster=='' || view.ae.Relationship=='' || view.ae.Dependency==''">
                    <option disabled value="">please select</option>
                    <option v-for="e in elements()">{{e}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Relationship&nbsp;Configuration:</td>
                <td>
                  <textarea id="aeConfiguration3"  rows=2
                    v-model="view.ae.Configuration3"
                    @change="configuration3Changed()"
                    @focus="configuration3Focus()"
                    @blur="configuration3Blur()"
                    :disabled="view.ae.Component=='' || view.ae.Cluster=='' || view.ae.Relationship==''"></textarea>
                </td>
              </tr>

            </table>
          </div>

          <div class="modal-footer">
            &nbsp;
            <button class="modal-default-button" v-if="!view.ae.New" v-on:click="deleteElement()">
              Delete Element <i class="fas fa-times-circle">
            </button>
            <button class="modal-default-button" v-if="!view.ae.New" v-on:click="updateElement()">
              Update Element <i class="fas fa-plus-circle">
            </button>
            <button class="modal-default-button" v-if="view.ae.New" v-on:click="createElement()">
              Create Element <i class="fas fa-plus-circle">
            </button>
          </div>
        </div>
      </div>`
  }
)
