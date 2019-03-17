Vue.component(
  'solutionElementEditor',
  {
    props: ['model', 'view'],
    methods: {
      configuration1Focus: function() {
        document.getElementById("seConfiguration1").setAttribute("rows", 10)
      },
      configuration1Blur: function() {
        document.getElementById("seConfiguration1").setAttribute("rows", 2)
      },
      configuration2Focus: function() {
        document.getElementById("seConfiguration2").setAttribute("rows", 10)
      },
      configuration2Blur: function() {
        document.getElementById("seConfiguration2").setAttribute("rows", 2)
      },
      configuration3Focus: function() {
        document.getElementById("seConfiguration3").setAttribute("rows", 10)
      },
      configuration3Blur: function() {
        document.getElementById("seConfiguration3").setAttribute("rows", 2)
      },
      // clusters creates a list of available cluster versions
      clusters: function() {
        var result = []
        for (var index in this.model.SolElement.Clusters) {
          version = this.model.SolElement.Clusters[index].Version
          result.push(version)
        }
        return result
      },
      // relationships
      relationships: function() {
        // check if the name and version and relationship have been defined
        if (this.view.se.Component == "" || this.view.se.Cluster == "") {
          return []
        }
        // match relationships
        cluster = this.model.SolElement.Clusters[this.view.se.Cluster]
        var result = []
        for (var index in cluster.Relationships) {
          relationship = cluster.Relationships[index]
          result.push(relationship.Relationship)
        }
        return result
      },
      // clusterChanged handles changes of the cluster information
      clusterChanged: function() {
        // determine cluster
        cluster = this.model.SolElement.Clusters[this.view.se.Cluster]

        // determine instances
        initial  = 0
        inactive = 0
        active   = 0
        failure  = 0
        other    = 0

        for (var index in cluster.Instances) {
          instance = cluster.Instances[index]
          switch (instance.State) {
            case "initial":
              initial++
              break;
            case "inactive":
              inactive++
              break;
            case "active":
              active++
              break;
            case "failure":
              failure++
              break;
            default:
              other++
          }
        }
        instances = "initial: " + initial + " inactive: " + inactive + " active: " + active + " failure:" + failure + " other: " + other

        // update the fields
        this.view.se.Target         = (cluster.Target != "" ? cluster.Target : "active")
        this.view.se.State          = (cluster.State  != "" ? cluster.State  : "active")
        this.view.se.Min            = (cluster.Min    != "" ? cluster.Min    : 1)
        this.view.se.Max            = (cluster.Max    != "" ? cluster.Max    : 1)
        this.view.se.Size           = (cluster.Size   != "" ? cluster.Size   : 1)
        this.view.se.Instances      = instances
        this.view.se.Configuration2 = cluster.Configuration
        this.view.se.Relationship   = ""
        this.view.se.Dependency     = ""
        this.view.se.RelElement     = ""
        this.view.se.Configuration3 = ""

        this.$forceUpdate()
      },
      // stateChanged  handles changes of the state information
      stateChanged: function() {
        this.model.SolElement.Clusters[this.view.se.Cluster].Status = this.view.se.Status
      },
      // minChanged  handles changes of the min information
      minChanged: function() {
        this.model.SolElement.Clusters[this.view.se.Cluster].Min = this.view.se.Min
      },
      // maxChanged  handles changes of the max information
      maxChanged: function() {
        this.model.SolElement.Clusters[this.view.se.Cluster].Max = this.view.se.Max
      },
      // sizeChanged  handles changes of the size information
      sizeChanged: function() {
        this.model.SolElement.Clusters[this.view.se.Cluster].Size = this.view.se.Size
      },
      // relationshipChanged handles changes of the relationship information
      relationshipChanged: function() {
        // determine cluster and relationship
        cluster = this.model.SolElement.Clusters[this.view.se.Cluster]
        relationship = cluster.Relationships[this.view.se.Relationship]

        // update the fields
        this.view.se.Dependency=     relationship.Dependency
        this.view.se.RelElement=     relationship.Element
        this.view.se.Configuration3= relationship.Configuration

        this.$forceUpdate()
      },

      // cancelDialog closes the wizard dialog
      cancelDialog: function() {
        this.model.SolElement = null
        this.view.se.New       = false
      },
      // updateElement updates an existing element of the architecture
      showTasks: function() {
        this.view.automation.solution = this.view.solution
        this.view.automation.element  = this.model.SolElement.Element
        this.view.automation.cluster  = this.view.se.Cluster
        this.view.automation.instance = ""

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)

        this.view.nav = "Automation"
      },
      // updateElement updates an existing element of the architecture
      updateElement: function() {
        Vue.delete(this.model.Architecture.Elements, this.view.se.Element)
        Vue.set(this.model.Architecture.Elements, this.view.se.Element, this.model.SolElement)
        this.model.SolElement = null
      }
    },
    template: `
      <div class="modal-mask" v-if="model.SolElement != null" @keyup.esc="cancelDialog()" tabindex="0">
        <div class="modal-container" @keyup.esc="cancelDialog()" tabindex="0">
          <div class="modal-header">
            <h3>Solution Element Viewer</h3>
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
                  <input type="text" v-model="view.se.Element" disabled></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Component:</td>
                <td>
                  <input type="text" v-model="view.se.Component" disabled></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Element&nbsp;Configuration:</td>
                <td>
                  <textarea id="seConfiguration1" rows=2
                    v-model="view.se.Configuration1"
                    @keyup.esc="configuration1Blur()"
                    @focus="configuration1Focus()"
                    @blur="configuration1Blur()"
                    :disabled="view.se.Element==''"></textarea>
                </td>
              </tr>

              <tr>
                <td><strong>Cluster:</strong></td>
                <td>
                  <select v-model="view.se.Cluster" v-on:change="clusterChanged()" :disabled="view.se.Component==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in clusters()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Target State:</td>
                <td>
                  <select v-model="view.se.Target" v-on:change="stateChanged()" disabled>
                    <option disabled>{{view.se.Target}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;State:</td>
                <td>
                  <select v-model="view.se.State" v-on:change="stateChanged()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''">
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
                  <input class="long" type="number" min="0" v-model="view.se.Min" v-on:change="minChanged()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Max:</td>
                <td>
                  <input class="long" type="number" min="0" v-model="view.se.Max" v-on:change="maxChanged()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Size:</td>
                <td>
                  <input class="long" type="number" min="0" v-model="view.se.Size" v-on:change="sizeChanged()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''"></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Instances:</td>
                <td>
                  <input class="long" v-model="view.se.Instances" disabled></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Cluster&nbsp;Configuration:</td>
                <td>
                  <textarea id="seConfiguration2" rows=2
                    v-model="view.se.Configuration2"
                    @keyup.esc="configuration2Blur()"
                    @focus="configuration2Focus()"
                    @blur="configuration2Blur()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''"></textarea>
                </td>
              </tr>

              <tr>
                <td><strong>Relationship:</strong></td>
                <td>
                  <select v-model="view.se.Relationship" v-on:change="relationshipChanged()"
                    :disabled="view.se.Component=='' || view.se.Cluster==''">
                    <option disabled value="">please select</option>
                    <option v-for="c in relationships()">{{c}}</option>
                  </select>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Dependency:</td>
                <td>
                  <input type="text" v-model="view.se.Dependency" disabled></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Element:</td>
                <td>
                  <input type="text" v-model="view.se.RelElement" disabled></input>
                </td>
              </tr>
              <tr>
                <td>&nbsp;Relationship&nbsp;Configuration:</td>
                <td>
                  <textarea id="seConfiguration3"  rows=2
                    v-model="view.se.Configuration3"
                    @keyup.esc="configuration3Blur()"
                    @focus="configuration3Focus()"
                    @blur="configuration3Blur()"
                    :disabled="view.se.Component=='' || view.se.Cluster=='' || view.se.Relationship==''"></textarea>
                </td>
              </tr>

            </table>
          </div>

          <div class="modal-footer">
            &nbsp;
            <button class="modal-default-button" v-on:click="showTasks()">
              Show Tasks <i class="fas fa-cogs">
            </button>
            <button class="modal-default-button" v-on:click="updateElement()">
              Update Element <i class="fas fa-plus-circle">
            </button>
          </div>
        </div>
      </div>`
  }
)
