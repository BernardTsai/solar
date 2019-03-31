Vue.component(
  'solEditor',
  {
    props: ['model', 'view','element'],
    methods: {
      // showStatus display status
      showStatus: function() {
        document.getElementById("solutionStatus").style.display        = "block"
        document.getElementById("solutionConfiguration").style.display = "none"
      },
      // showConfiguration display configuration
      showConfiguration: function() {
        document.getElementById("solutionStatus").style.display        = "none"
        document.getElementById("solutionConfiguration").style.display = "block"
      },
      // showTasks switches to the corresponding automation tasks
      showTasks: function() {
        model.Tasks = []
        model.Trace = null

        view.automation.solution = view.solution
        view.automation.element  = this.element.Element
        view.automation.cluster  = ""
        view.automation.instance = ""

        loadTasks(
          view.domain,
          view.automation.solution,
          view.automation.element,
          view.automation.cluster,
          view.automation.instance)

        view.nav = "Automation"
      },
      // editConfiguration opens editor for editing a configuration
      editConfiguration: function(rel) {
        if (rel != "") {
          this.view.se.ConfTitle     = "Configuration of relationship: " + rel.Relationship
          this.view.se.Configuration = rel.Configuration
          this.view.se.Dependency    = rel.Dependency
        } else {
          this.view.se.ConfTitle     = "Configuration of element: "
          this.view.se.Configuration = this.element.Configuration
          this.view.se.Dependency    = ""
        }
        this.$forceUpdate()
      },
      // discardConfiguration closes the configuration editor
      discardConfiguration: function() {
        this.view.se.Configuration = null
        this.view.se.Dependency    = ""

        this.$forceUpdate()
      },
      // updateClusterSize starts a task to update the size of a cluster
      updateClusterSize: function(cluster, min, max, size) {
        updateCluster(this.view.domain, this.view.solution, this.model.SolElement.Element, cluster.Version, cluster.State, min, max, size)
        this.$emit('refresh')
      },
      // updateClusterState starts a task to update the state of a cluster
      updateClusterState: function(cluster, state) {
        updateCluster(this.view.domain, this.view.solution, this.model.SolElement.Element, cluster.Version, state, cluster.Min, cluster.Max, cluster.Size)
        this.$emit('refresh')
      },
      // updateInstanceState starts a task to update the state of an instance
      updateInstanceState: function(cluster, instance, state) {
        updateInstance(this.view.domain, this.view.solution, this.model.SolElement.Element, cluster.Version, instance.UUID, state)
        this.$emit('refresh')
      }
    },
    template: `
      <div class="solEditor">
        <div class="header">
          <h3>Element: {{element.Element}} ({{element.Component}})</h3>
          <button class="modal-default-button" @click="showTasks()">
            Tasks <i class="fas fa-cogs">
          </button>
          <button class="modal-default-button" @click="showStatus()">
            Status <i class="fas fa-heartbeat">
          </button>
          <button class="modal-default-button" @click="showConfiguration()">
            Configuration <i class="fas fa-file-alt">
          </button>
        </div>

        <div id="solutionConfiguration" class="configuration">
          <table style="width: 100%">
            <col width="10*">
            <col width="990*">
            <tr>
              <td><strong>Element:</strong></td>
              <td>
                <input type="text" readonly
                  v-model="element.Element">
                </input>
              </td>
            </tr>
            <tr>
              <td>&nbsp;Component:</td>
              <td>
                <input type="text" readonly
                  v-model="element.Component">
                </input>
              </td>
            </tr>
            <tr>
              <td>&nbsp;Configuration:</td>
              <td>
                <textarea id="configuration" rows=5 readonly
                  @click="editConfiguration('')"
                  v-model="element.Configuration">
                </textarea>
              </td>
            </tr>

            <tr>
              <td><strong>Clusters:</strong></td>
              <td>
                <select
                  v-model="view.se.Cluster">
                  <option value="">please select</option>
                  <option v-for="cluster in element.Clusters">{{cluster.Version}}</option>
                </select>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Target State:</td>
              <td>
                <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Target"/>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Current State:</td>
              <td>
                <input type="text" readonly v-model="element.Clusters[view.se.Cluster].State"/>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Minimum Size:</td>
              <td>
                <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Min"/>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Maximum Size:</td>
              <td>
                <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Max"/>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Current Size:</td>
              <td>
                <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Size"/>
              </td>
            </tr>
            <tr v-if="view.se.Cluster != ''">
              <td>&nbsp;Configuration:</td>
              <td>
                <textarea id="configuration2" rows=5 readonly
                  @click="editConfiguration('')"
                  v-model="element.Clusters[view.se.Cluster].Configuration">
                </textarea>
              </td>
            </tr>

            <tr v-if="view.se.Cluster != ''">
              <td><strong>Relationships:</strong></td>
              <td>

                <table class="relationships">
                  <thead>
                    <tr>
                      <th>Relationship</th>
                      <th>Dependency</th>
                      <th>Element</th>
                      <th>Configuration</th>
                  </thead>
                  <tbody>
                    <tr v-for="rel in element.Clusters[view.se.Cluster].Relationships">
                      <td>
                        <input type="text" readonly v-model="rel.Relationship"/>
                      </td>
                      <td>
                        <input type="text" readonly v-model="rel.Dependency"/>
                      </td>
                      <td>
                        <input type="text" readonly v-model="rel.Element"/>
                      </td>
                      <td @click="editConfiguration(rel)">
                        <i class="fas fa-edit"></i>
                      </td>
                    </tr>
                  </tbody>
                </table>

              </td>
            </tr>
          </table>
        </div>

        <div id="solutionStatus" class="status">
          <table class="clusters">

            <tr v-for="cluster in element.Clusters">
              <td>
                <table class="cluster">
                  <tr>
                    <th class="label">Cluster: {{cluster.Version}}</th>
                    <th style="text-align:center;">Minimum</th>
                    <th style="text-align:center;">Maximum</th>
                    <th style="text-align:center;">Size</th>
                  </tr>
                  <tr>
                    <td class="label">&nbsp;&nbsp;&nbsp;&nbsp;current&nbsp;size</td>
                    <td class="number">{{cluster.Min}}</td>
                    <td class="number">{{cluster.Max}}</td>
                    <td class="number">{{cluster.Size}}</td>
                  </tr>
                  <tr>
                    <td class="label">&nbsp;&nbsp;&nbsp;&nbsp;target&nbsp;size</td>
                    <td class="number">
                      <input type="text" v-model="cluster.NewMin"/>
                    </td>
                    <td class="number">
                      <input type="text" v-model="cluster.NewMax"/>
                    </td>
                    <td class="number">
                      <input type="text" v-model="cluster.NewSize"/>
                    </td>
                    <td>
                      <button class="modal-default-button"
                        @click="updateClusterSize(cluster, parseInt(cluster.NewMin), parseInt(cluster.NewMax), parseInt(cluster.NewSize))">Update <i class="fas fa-redo"/>
                      </button>
                    </td>
                  </tr>

                  <tr>
                    <th>&nbsp;&nbsp;Cluster</th>
                    <td class="state initial"
                      :class="{selected: cluster.State=='initial'}"
                      @click="updateClusterState(cluster,'initial')">inital</td>
                    <td class="state inactive"
                      :class="{selected: cluster.State=='inactive'}"
                      @click="updateClusterState(cluster,'inactive')">inactive</td>
                    <td class="state active"
                      :class="{selected: cluster.State=='active'}"
                      @click="updateClusterState(cluster,'active')">active</td>
                  </tr>
                  <tr v-for="instance in _.orderBy(cluster.Instances, 'UUID')">
                    <td>&nbsp;&nbsp;&bullet; {{instance.UUID}}</td>
                    <td class="state initial"
                      :class="{selected: instance.State=='initial'}"
                      @click="updateInstanceState(cluster,instance,'initial')">inital</td>
                    <td class="state inactive"
                      :class="{selected: instance.State=='inactive'}"
                      @click="updateInstanceState(cluster,instance,'inactive')">inactive</td>
                    <td class="state active"
                      :class="{selected: instance.State=='active'}"
                      @click="updateInstanceState(cluster,instance,'active')">active</td>
                    <td class="state failure"
                      :class="{selected: instance.State=='failure'}"
                      @click="updateInstanceState(cluster,instance,'failure')">failure</td>
                  </tr>

                </table>
              </td>
            </tr>

          </table>
        </div>



        <div class="configurationEditor" v-if="view.se.Configuration != null">
          <div class="modal">
            <h3>{{view.se.ConfTitle}}</h3>
            <textarea v-focus @keyup.esc="discardConfiguration()" v-model="view.se.Configuration"></textarea>
            <button class="modal-default-button" @click="discardConfiguration()">
              Cancel <i class="fas fa-times-circle">
            </button>

          </div>
        </div>

      </div>`
  }
)
