Vue.component(
  'architecture',
  {
    props: ['model', 'view'],
    methods: {
      // create creates a new solution
      create: function() {
        // ask for name of new solution and add a new architecture
        name = prompt("Name of the new solution:")
        if (name != null && name != "" && name != "null") {
          architecture = {
            Architecture:  name,
            Version:       "V0.0.0",
            Configuration: "",
            Elements:      {}
          }

          saveArchitecture(this.view.domain, architecture)
          .then(() => { loadArchitectures(this.view.domain) })
          .then(() => {
            this.view.architecture = name + " - V0.0.0"
            this.selectArchitecture()
          })
        }
      },
      // remove deletes the current architecture
      remove: function() {
        deleteArchitecture(this.view.domain, this.model.Architecture)
        .then(() => { loadArchitectures(this.view.domain) })
        .then(() => {
          this.view.architecture = ""
          this.selectArchitecture()
        })
      },
      // edit allows to specify configuration information
      edit: function() {
        // hide element editor
        this.model.ArchElement = null
        // hide graph
        this.view.showGraph    = false
      },
      // update updates the architecture
      update: function() {
        architecture = this.model.Architecture

        // cleanup strings which should be numbers
        Object.values(architecture.Elements).forEach((element) => {
          Object.values(element.Clusters).forEach((cluster) => {
            cluster.Min  = parseInt(cluster.Min)
            cluster.Max  = parseInt(cluster.Max)
            cluster.Size = parseInt(cluster.Size)
          })
        })

        saveArchitecture(this.view.domain, architecture)
        .then(() => {
          this.view.architecture = architecture.Architecture + " - " + architecture.Version
          this.selectArchitecture()
        })
      },
      // duplicate creates a copy of the architecture
      duplicate: function() {
        // ask for version of new architecture and create a copy of the existing architecture
        version = prompt("Version of the architecture:")
        if (version != null && version != "" && version != "null") {
          architecture = duplicateArchitecture(this.model.Architecture, version)

          // cleanup strings which should be numbers
          Object.values(architecture.Elements).forEach((element) => {
            Object.values(element.Clusters).forEach((cluster) => {
              cluster.Min  = parseInt(cluster.Min)
              cluster.Max  = parseInt(cluster.Max)
              cluster.Size = parseInt(cluster.Size)
            })
          })

          // save the new architecture and switch to it
          saveArchitecture(this.view.domain, architecture)
          .then(() => { loadArchitectures(this.view.domain) })
          .then(() => {
            this.view.architecture = architecture.Architecture + " - " + version
            this.selectArchitecture()
          })
        }
      },
      // deploy deploys the architecture blueprint
      deploy: function() {
        architecture = this.model.Architecture

        deployArchitecture(this.view.domain, architecture)
        .then(() => loadSolutions(this.view.domain) )
        .then(() => navSolution() )
      },
      // viewNode displays a nodes element in the editor
      viewNode: function(node) {
        // initialise the architecture element of the model
        this.model.ArchElement = node.Element

        // add missing clusters
        Object.values(this.model.Catalog).forEach((c) => {
          if (c.Component == node.Element.Component && !node.Element.Clusters[c.Version]) {
            node.Element.Clusters[c.Version] = {
              Version:       c.Version,
              State:         "initial",
              Min:           0,
              Max:           0,
              Size:          0,
              Configuration: "",
              Relationships: {}
            }
          }
        })

      },
      // hidelement hides the editor
      hideElement: function() {
        // reset the architecture element of the model
        this.model.ArchElement = null

        // show graph
        this.view.showGraph = true

        this.$forceUpdate()
      },
      // addElement adds an element to the architecture based on a component
      addElement: function(component) {
        if (this.view.architecture != '') {
          // ask for name of new element and add it with initial values
          name = prompt("Name of the new '" + component + "' element:")
          if (name != null && name != "" && name != "null") {
            // create new element
            element = {
              Element:       name,
              Component:     component,
              Configuration: "",
              Clusters:      {}
            }

            // add element to architecture
            Vue.set(this.model.Architecture.Elements, name, element )

            // add available clusters
            Object.values(this.model.Catalog).forEach((c) => {
              if (c.Component == component) {
                element.Clusters[c.Version] = {
                  Version:       c.Version,
                  State:         "active",
                  Min:           1,
                  Max:           1,
                  Size:          1,
                  Configuration: "",
                  Relationships: {}
                }
              }
            })

            this.$forceUpdate()

            this.model.ArchElement = this.model.Architecture.Elements[name]
          }
        }
      },
      // selectArchitecture pick a specific version of an architecture
       selectArchitecture: function() {
         // load architecture
        if (this.view.domain != "" && this.view.architecture != ""){
          this.view.solution = getName(this.view.architecture)
          this.view.version  = getVersion(this.view.architecture)

          loadArchitecture(this.view.domain, this.view.architecture)
        } else {
          this.model.Architecture = null
          this.model.Graph        = null
        }
        // reset element
        this.model.ArchElement = null
      },
      // graph creates the architecture graph
      graph: function() {
        return new ArchitectureGraph(this.model, this.view, this.view.domain, this.view.solution, this.view.version)
      }
    },
    computed: {
      // components provides a sorted list of component names
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

        // return the sorted result
        result.sort()

        return result
      }
    },
    template: `
      <div id="architecture" v-if="view.nav=='Architecture'">

        <div id="selector">
          <div id="architecture-selector">
            <strong>Architecture:</strong>
            <select id="architectureSelector" v-model="view.architecture" @change="selectArchitecture">
              <option selected value="">Please select one</option>
              <option v-for="architecture in model.Architectures">{{architecture}}</option>
            </select>
          </div>

          <div class="buttons">
            <button class="action" v-if="view.architecture==''" @click="create()" title="Create a new architecture">
              Create <i class="fas fa-plus-circle">
            </button>
            <button class="action" v-if="view.architecture!=''" @click="edit()" title="Edit architecture information">
              Edit <i class="fas fa-edit">
            </button>
            <button class="action" v-if="view.architecture!=''" @click="update()" title="Save the architecture">
              Update <i class="fas fa-cloud-upload-alt">
            </button>
            <button class="action" v-if="view.architecture!=''" @click="deploy()" title="Deploy the architecture">
              Deploy <i class="fas fa-play-circle">
            </button>
            <button class="action" v-if="view.architecture!=''" @click="duplicate()" title="Create a copy of the architecture">
              Duplicate <i class="fas fa-copy">
            </button>
            <button class="action" v-if="view.architecture!=''" @click="remove()" title="Delete the architecture">
              Delete <i class="fas fa-minus-circle">
            </button>
          </div>
        </div>

        <div id="comps" v-if="model.Architecture">
          <div class="header">
            <h3>Catalog:</h3>
          </div>

          <table class="components">
            <thead>
              <tr>
                <th>Component</th>
                <th @click="hideElement()">
                  <i class="fas fa-eye-slash"></i>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="component in components">
                <td>{{component}}</td>
                <td  @click="addElement(component)">
                  <i class="fas fa-cube"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <architectureEditor v-if="model.Architecture && !view.showGraph" :model="model" :view="view"/>

        <div id="container" v-if="model.Architecture && view.showGraph">
          <graph :model="model" :view="view" :graph="graph()" @node-selected="viewNode"/>
        </div>

        <elementEditor v-if="model.ArchElement" :model="model" :view="view" :element="model.ArchElement"/>

      </div>`
  }
)
